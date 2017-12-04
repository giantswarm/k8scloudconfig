/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	v1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	scheme "github.com/giantswarm/apiextensions/pkg/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// IngressConfigsGetter has a method to return a IngressConfigInterface.
// A group's client should implement this interface.
type IngressConfigsGetter interface {
	IngressConfigs(namespace string) IngressConfigInterface
}

// IngressConfigInterface has methods to work with IngressConfig resources.
type IngressConfigInterface interface {
	Create(*v1alpha1.IngressConfig) (*v1alpha1.IngressConfig, error)
	Update(*v1alpha1.IngressConfig) (*v1alpha1.IngressConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.IngressConfig, error)
	List(opts v1.ListOptions) (*v1alpha1.IngressConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.IngressConfig, err error)
	IngressConfigExpansion
}

// ingressConfigs implements IngressConfigInterface
type ingressConfigs struct {
	client rest.Interface
	ns     string
}

// newIngressConfigs returns a IngressConfigs
func newIngressConfigs(c *CoreV1alpha1Client, namespace string) *ingressConfigs {
	return &ingressConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the ingressConfig, and returns the corresponding ingressConfig object, and an error if there is any.
func (c *ingressConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.IngressConfig, err error) {
	result = &v1alpha1.IngressConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("ingressconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of IngressConfigs that match those selectors.
func (c *ingressConfigs) List(opts v1.ListOptions) (result *v1alpha1.IngressConfigList, err error) {
	result = &v1alpha1.IngressConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("ingressconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested ingressConfigs.
func (c *ingressConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("ingressconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a ingressConfig and creates it.  Returns the server's representation of the ingressConfig, and an error, if there is any.
func (c *ingressConfigs) Create(ingressConfig *v1alpha1.IngressConfig) (result *v1alpha1.IngressConfig, err error) {
	result = &v1alpha1.IngressConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("ingressconfigs").
		Body(ingressConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a ingressConfig and updates it. Returns the server's representation of the ingressConfig, and an error, if there is any.
func (c *ingressConfigs) Update(ingressConfig *v1alpha1.IngressConfig) (result *v1alpha1.IngressConfig, err error) {
	result = &v1alpha1.IngressConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("ingressconfigs").
		Name(ingressConfig.Name).
		Body(ingressConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the ingressConfig and deletes it. Returns an error if one occurs.
func (c *ingressConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("ingressconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *ingressConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("ingressconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched ingressConfig.
func (c *ingressConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.IngressConfig, err error) {
	result = &v1alpha1.IngressConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("ingressconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
