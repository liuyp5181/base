package proxy

import (
	"context"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// Proxy is a dynamic gRPC client that performs reflection
type Proxy struct {
	cc        *grpc.ClientConn
	reflector Reflector
	stub      Stub
}

// NewConnect opens a connection to target.
func NewConnect(ctx context.Context, target string) (*Proxy, error) {
	p := &Proxy{}
	cc, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	p.cc = cc
	rc := grpcreflect.NewClient(ctx, rpb.NewServerReflectionClient(p.cc))
	p.reflector = NewReflector(rc)
	p.stub = NewStub(grpcdynamic.NewStub(p.cc))
	return p, nil
}

// NewClient opens a connection to target.
func NewClient(ctx context.Context, cc *grpc.ClientConn) *Proxy {
	p := &Proxy{}
	p.cc = cc
	rc := grpcreflect.NewClient(ctx, rpb.NewServerReflectionClient(p.cc))
	p.reflector = NewReflector(rc)
	p.stub = NewStub(grpcdynamic.NewStub(p.cc))
	return p
}

// CloseConn closes the underlying connection
func (p *Proxy) CloseConn() error {
	return p.cc.Close()
}

// Call performs the gRPC call after doing reflection to obtain type information
func (p *Proxy) Call(ctx context.Context,
	serviceName, methodName string,
	message []byte,
	opts ...grpc.CallOption,
) ([]byte, error) {

	invocation, err := p.reflector.CreateInvocation(ctx, serviceName, methodName, message)
	if err != nil {
		return nil, err
	}

	outputMsg, err := p.stub.InvokeRPC(ctx, invocation, opts...)
	if err != nil {
		return nil, err
	}
	m, err := outputMsg.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal output JSON")
	}
	return m, err
}
