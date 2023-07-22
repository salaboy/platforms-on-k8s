package v1alpha1

import (
	"context"

	"github.com/salaboy/platforms-on-k8s/conference-admin/admin-go/api/types/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ConferenceAdminV1Alpha1Interface interface {
	Environments(namespace string) EnvironmentInterface
}

type ConferenceAdminV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ConferenceAdminV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ConferenceAdminV1Alpha1Client{restClient: client}, nil
}

func (c *ConferenceAdminV1Alpha1Client) Environments(namespace string) EnvironmentInterface {
	return &envClient{
		restClient: c.restClient,
		ns:         namespace,
		ctx:        context.TODO(),
	}
}
