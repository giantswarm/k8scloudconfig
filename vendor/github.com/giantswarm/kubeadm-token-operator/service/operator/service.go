package operator

import (
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/kubeadmtokentpr"
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/giantswarm/operatorkit/client/k8s"
	"github.com/giantswarm/operatorkit/operator"
	"github.com/giantswarm/operatorkit/tpr"
	"github.com/juju/errgo"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/tools/cache"

	tokenutil "github.com/giantswarm/kubeadm-token-operator/token"
)

const (
	ClusterIDLabel  = "clusterID"
	KubeadmTokenKey = "kubeadmToken"

	deleteSecretMaxElapsedTime = 30 * time.Second
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	// Settings.
	Namespace string
	Service   string
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	var err error

	var k8sClient kubernetes.Interface
	{
		config := k8s.DefaultConfig()
		k8sClient, err = k8s.NewClient(config)
		if err != nil {
			panic(err)
		}
	}

	var newLogger micrologger.Logger
	{
		config := micrologger.DefaultConfig()
		newLogger, err = micrologger.New(config)
		if err != nil {
			panic(err)
		}
	}

	return Config{
		// Dependencies.
		K8sClient: k8sClient,
		Logger:    newLogger,
	}
}

// New creates a new configured service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.K8sClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}

	var err error
	var newTPR *tpr.TPR
	{
		tprConfig := tpr.DefaultConfig()

		tprConfig.K8sClient = config.K8sClient
		tprConfig.Logger = config.Logger

		tprConfig.Description = kubeadmtokentpr.Description
		tprConfig.Name = kubeadmtokentpr.Name
		tprConfig.Version = kubeadmtokentpr.VersionV1

		newTPR, err = tpr.New(tprConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newService := &Service{
		// Dependencies.
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		// Internals
		bootOnce: sync.Once{},
		mutex:    sync.Mutex{},
		tpr:      newTPR,
	}

	return newService, nil
}

// Service implements the service.
type Service struct {
	// Dependencies.
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	// Internals.
	ipsPorts  map[string]int
	bootOnce  sync.Once
	mutex     sync.Mutex
	namespace string
	service   string
	tpr       *tpr.TPR
}

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		err := s.tpr.CreateAndWait()
		if tpr.IsAlreadyExists(err) {
			s.logger.Log("debug", "third party resource already exists")
		} else if err != nil {
			s.logger.Log("error", errgo.Details(err))
			return
		}

		s.logger.Log("debug", "starting list/watch")

		newResourceEventHandler := &cache.ResourceEventHandlerFuncs{
			AddFunc:    s.addFunc,
			DeleteFunc: s.deleteFunc,
		}
		newZeroObjectFactory := &tpr.ZeroObjectFactoryFuncs{
			NewObjectFunc:     func() runtime.Object { return &kubeadmtokentpr.CustomObject{} },
			NewObjectListFunc: func() runtime.Object { return &kubeadmtokentpr.List{} },
		}

		s.tpr.NewInformer(newResourceEventHandler, newZeroObjectFactory).Run(nil)

	})
}

func secretName(obj kubeadmtokentpr.CustomObject) string {
	return fmt.Sprintf("%s-token", obj.Spec.ClusterID)
}

func (s *Service) addFunc(obj interface{}) {
	// We lock the addFunc/deleteFunc to make sure only one addFunc/deleteFunc is
	// executed at a time. addFunc/deleteFunc is not thread safe. This is
	// important because the source of truth for the ingress-operator are
	// Kubernetes resources. In case we would run the operator logic in parallel,
	// we would run into race conditions.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := operator.ProcessCreate(obj, s); err != nil {
		s.logger.Log("error", errgo.Details(err))
	}
}

func (s *Service) deleteSecret(obj kubeadmtokentpr.CustomObject) error {
	err := s.k8sClient.Core().Secrets(v1.NamespaceDefault).Delete(secretName(obj), &v1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}

func (s *Service) deleteSecretAndWait(obj kubeadmtokentpr.CustomObject) error {
	initBackoff := backoff.NewExponentialBackOff()
	initBackoff.MaxElapsedTime = deleteSecretMaxElapsedTime

	operation := func() error {
		err := s.deleteSecret(obj)
		if err != nil {
			return microerror.MaskAny(err)
		}

		return nil
	}

	notify := func(reason error, interval time.Duration) {
		s.logger.Log("debug", "failed to delete secret: %s", errgo.Details(reason))
	}

	return microerror.MaskAny(backoff.RetryNotify(operation, initBackoff, notify))
}

func (s *Service) deleteFunc(obj interface{}) {
	// We lock the addFunc/deleteFunc to make sure only one addFunc/deleteFunc is
	// executed at a time. addFunc/deleteFunc is not thread safe. This is
	// important because the source of truth for the ingress-operator are
	// Kubernetes resources. In case we would run the operator logic in parallel,
	// we would run into race conditions.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := operator.ProcessDelete(obj, s); err != nil {
		s.logger.Log("error", errgo.Details(err))
	}
}

func (s *Service) GetCurrentState(obj interface{}) (interface{}, error) {
	customObject, ok := obj.(*kubeadmtokentpr.CustomObject)
	if !ok {
		return nil, microerror.MaskAnyf(wrongTypeError, "expected '%T', got '%T'", &kubeadmtokentpr.CustomObject{}, obj)
	}

	var cState CurrentState

	secret, err := s.k8sClient.CoreV1().Secrets(v1.NamespaceDefault).Get(secretName(*customObject))
	if errors.IsNotFound(err) {
		return cState, nil
	} else if err != nil {
		return nil, microerror.MaskAny(err)
	}

	cState.Secret = secret

	return cState, nil
}

func (s *Service) GetDesiredState(obj interface{}) (interface{}, error) {
	customObject, ok := obj.(*kubeadmtokentpr.CustomObject)
	if !ok {
		return nil, microerror.MaskAnyf(wrongTypeError, "expected '%T', got '%T'", &kubeadmtokentpr.CustomObject{}, obj)
	}

	token, err := tokenutil.GenerateToken()
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	var dState DesiredState

	dState.Secret = &v1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name: secretName(*customObject),
			Labels: map[string]string{
				ClusterIDLabel: customObject.Spec.ClusterID,
			},
		},
		StringData: map[string]string{
			KubeadmTokenKey: token,
		},
	}

	return dState, nil
}

func (s *Service) GetCreateState(obj, currentState, desiredState interface{}) (interface{}, error) {
	cState, ok := currentState.(CurrentState)
	if !ok {
		return nil, microerror.MaskAnyf(wrongTypeError, "expected '%T', got '%T'", CurrentState{}, currentState)
	}
	dState, ok := desiredState.(DesiredState)
	if !ok {
		return nil, microerror.MaskAnyf(wrongTypeError, "expected '%T', got '%T'", DesiredState{}, desiredState)
	}

	var createState ActionState

	if cState.Secret == nil {
		createState.Secret = dState.Secret
	}

	return createState, nil
}

func (s *Service) GetDeleteState(obj, currentState, desiredState interface{}) (interface{}, error) {
	dState, ok := desiredState.(DesiredState)
	if !ok {
		return nil, microerror.MaskAnyf(wrongTypeError, "expected '%T', got '%T'", DesiredState{}, desiredState)
	}

	var deleteState ActionState

	if dState.Secret != nil {
		deleteState.Secret = nil
	}

	return deleteState, nil
}

func (s *Service) ProcessCreateState(obj, createState interface{}) error {
	customObject, ok := obj.(*kubeadmtokentpr.CustomObject)
	if !ok {
		return microerror.MaskAnyf(wrongTypeError, "expected '%T', got '%T'", &kubeadmtokentpr.CustomObject{}, obj)
	}
	cState, ok := createState.(ActionState)
	if !ok {
		return microerror.MaskAnyf(wrongTypeError, "expected '%T', got '%T'", ActionState{}, createState)
	}

	s.logger.Log("debug", "process create state", "cluster", customObject.Spec.ClusterID)

	if cState.Secret == nil {
		return nil
	}

	_, err := s.k8sClient.Core().Secrets(v1.NamespaceDefault).Create(cState.Secret)
	if err != nil {
		return microerror.MaskAny(err)
	}

	s.logger.Log("info", "secret with token created", "cluster", customObject.Spec.ClusterID)

	return nil
}

func (s *Service) ProcessDeleteState(obj, deleteState interface{}) error {
	customObject, ok := obj.(*kubeadmtokentpr.CustomObject)
	if !ok {
		return microerror.MaskAnyf(wrongTypeError, "expected '%T', got '%T'", &kubeadmtokentpr.CustomObject{}, obj)
	}

	s.logger.Log("debug", "process create state", "cluster", customObject.Spec.ClusterID)

	if err := s.k8sClient.Core().Secrets(v1.NamespaceDefault).Delete(secretName(*customObject), &v1.DeleteOptions{}); err != nil {
		return microerror.MaskAny(err)
	}

	s.logger.Log("info", "secret with token deleted", "cluster", customObject.Spec.ClusterID)

	return nil
}
