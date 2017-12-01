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

// FlannelConfigsGetter has a method to return a FlannelConfigInterface.
// A group's client should implement this interface.
type FlannelConfigsGetter interface {
	FlannelConfigs(namespace string) FlannelConfigInterface
}

// FlannelConfigInterface has methods to work with FlannelConfig resources.
type FlannelConfigInterface interface {
	Create(*v1alpha1.FlannelConfig) (*v1alpha1.FlannelConfig, error)
	Update(*v1alpha1.FlannelConfig) (*v1alpha1.FlannelConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.FlannelConfig, error)
	List(opts v1.ListOptions) (*v1alpha1.FlannelConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.FlannelConfig, err error)
	FlannelConfigExpansion
}

// flannelConfigs implements FlannelConfigInterface
type flannelConfigs struct {
	client rest.Interface
	ns     string
}

// newFlannelConfigs returns a FlannelConfigs
func newFlannelConfigs(c *CoreV1alpha1Client, namespace string) *flannelConfigs {
	return &flannelConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the flannelConfig, and returns the corresponding flannelConfig object, and an error if there is any.
func (c *flannelConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.FlannelConfig, err error) {
	result = &v1alpha1.FlannelConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("flannelconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FlannelConfigs that match those selectors.
func (c *flannelConfigs) List(opts v1.ListOptions) (result *v1alpha1.FlannelConfigList, err error) {
	result = &v1alpha1.FlannelConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("flannelconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested flannelConfigs.
func (c *flannelConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("flannelconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a flannelConfig and creates it.  Returns the server's representation of the flannelConfig, and an error, if there is any.
func (c *flannelConfigs) Create(flannelConfig *v1alpha1.FlannelConfig) (result *v1alpha1.FlannelConfig, err error) {
	result = &v1alpha1.FlannelConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("flannelconfigs").
		Body(flannelConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a flannelConfig and updates it. Returns the server's representation of the flannelConfig, and an error, if there is any.
func (c *flannelConfigs) Update(flannelConfig *v1alpha1.FlannelConfig) (result *v1alpha1.FlannelConfig, err error) {
	result = &v1alpha1.FlannelConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("flannelconfigs").
		Name(flannelConfig.Name).
		Body(flannelConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the flannelConfig and deletes it. Returns an error if one occurs.
func (c *flannelConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("flannelconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *flannelConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("flannelconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched flannelConfig.
func (c *flannelConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.FlannelConfig, err error) {
	result = &v1alpha1.FlannelConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("flannelconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
