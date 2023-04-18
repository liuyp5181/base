package extend

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Extend struct {
	Ctx context.Context
	MD  *metadata.MD
}

func New() *Extend {
	return &Extend{
		Ctx: context.Background(),
		MD:  &metadata.MD{},
	}
}

func NewContext(ctx context.Context) *Extend {
	return &Extend{
		Ctx: ctx,
		MD:  &metadata.MD{},
	}
}

func (e *Extend) SetClient(key, val string) {
	e.Ctx = metadata.AppendToOutgoingContext(e.Ctx, key, val)
}

func (e *Extend) GetClient(key string) string {
	md, ok := metadata.FromIncomingContext(e.Ctx)
	if !ok {
		return ""
	}
	if len(md[key]) == 0 {
		return ""
	}
	return md[key][0]
}

func (e *Extend) SetServer(key, val string) {
	grpc.SetHeader(e.Ctx, metadata.MD{key: {val}})
}

func (e *Extend) LoadServer() grpc.CallOption {
	return grpc.Header(e.MD)
}

func (e *Extend) GetServer(key string) string {
	v := e.MD.Get(key)
	if len(v) == 0 {
		return ""
	}
	return v[0]
}
