package operator

import (
	"k8s.io/client-go/pkg/api/v1"
)

type ActionState struct {
	Secret *v1.Secret
}

type CurrentState struct {
	Secret *v1.Secret
}

type DesiredState struct {
	Secret *v1.Secret
}
