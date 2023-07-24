package v1alpha1

import (
	"context"

	"github.com/salaboy/platforms-on-k8s/conference-admin/admin-go/api/types/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type EnvironmentInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.EnvironmentList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Environment, error)
	Create(*v1alpha1.Environment) (*v1alpha1.Environment, error)
	Delete(name string, opts metav1.DeleteOptions) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type envClient struct {
	restClient rest.Interface
	ns         string
	ctx        context.Context
}

func (c *envClient) List(opts metav1.ListOptions) (*v1alpha1.EnvironmentList, error) {
	result := v1alpha1.EnvironmentList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("environments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *envClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.Environment, error) {
	result := v1alpha1.Environment{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("environments").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *envClient) Create(environment *v1alpha1.Environment) (*v1alpha1.Environment, error) {
	result := v1alpha1.Environment{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("environments").
		Body(environment).
		Do(c.ctx).
		Into(&result)

	return &result, err
}

func (c *envClient) Delete(name string, opts metav1.DeleteOptions) error {
	return c.restClient.
		Delete().
		Namespace(c.ns).
		Resource("environments").
		Name(name).
		Body(&opts).
		Do(c.ctx).
		Error()
}

func (c *envClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("environments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(c.ctx)
}
