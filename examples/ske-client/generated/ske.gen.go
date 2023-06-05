// Package ske provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/do87/stackit-client-generator version (devel) DO NOT EDIT.
package ske

import (
	"net/url"
	"strings"

	skeclient "github.com/do87/stackit-client-generator/examples/ske-client"
	"github.com/do87/stackit-client-generator/examples/ske-client/generated/cluster"
	"github.com/do87/stackit-client-generator/examples/ske-client/generated/credentials"
	"github.com/do87/stackit-client-generator/examples/ske-client/generated/operation"
	"github.com/do87/stackit-client-generator/examples/ske-client/generated/project"
	provideroptions "github.com/do87/stackit-client-generator/examples/ske-client/generated/provider-options"
)

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// list of connected client services
	Cluster         *cluster.Client
	Credentials     *credentials.Client
	Operation       *operation.Client
	ProviderOptions *provideroptions.Client
	Project         *project.Client

	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client skeclient.HttpRequestDoer
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

func NewRawClient(server string, opts ...ClientOption) (*Client, error) {
	// create a factory client
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}

	client.Cluster = cluster.NewRawClient(server, client.Client)
	client.Credentials = credentials.NewRawClient(server, client.Client)
	client.Operation = operation.NewRawClient(server, client.Client)
	client.ProviderOptions = provideroptions.NewRawClient(server, client.Client)
	client.Project = project.NewRawClient(server, client.Client)

	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer skeclient.HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponses builds on rawClientInterface to offer response payloads
type ClientWithResponses struct {
	Client *Client

	// list of connected client services
	Cluster         *cluster.ClientWithResponses
	Credentials     *credentials.ClientWithResponses
	Operation       *operation.ClientWithResponses
	ProviderOptions *provideroptions.ClientWithResponses
	Project         *project.ClientWithResponses
}

// NewClient creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClient(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewRawClient(server, opts...)
	if err != nil {
		return nil, err
	}

	cwr := &ClientWithResponses{Client: client}
	cwr.Cluster = cluster.NewClient(server, client.Client)
	cwr.Credentials = credentials.NewClient(server, client.Client)
	cwr.Operation = operation.NewClient(server, client.Client)
	cwr.ProviderOptions = provideroptions.NewClient(server, client.Client)
	cwr.Project = project.NewClient(server, client.Client)

	return cwr, nil
}
