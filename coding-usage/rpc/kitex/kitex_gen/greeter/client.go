// Code generated by Kitex v0.13.1. DO NOT EDIT.

package greeter

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	kitex_gen "github.com/luvx21/coding-go/coding-usage/rpc/kitex/kitex_gen"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	SayHello(ctx context.Context, Req *kitex_gen.HelloRequest, callOptions ...callopt.Option) (r *kitex_gen.HelloReply, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kGreeterClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kGreeterClient struct {
	*kClient
}

func (p *kGreeterClient) SayHello(ctx context.Context, Req *kitex_gen.HelloRequest, callOptions ...callopt.Option) (r *kitex_gen.HelloReply, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SayHello(ctx, Req)
}
