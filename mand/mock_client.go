package mand

import (
	context "context"

	grpc "google.golang.org/grpc"
)

type MockClientConn struct{}

func NewMockClientConn() *MockClientConn {
	return &MockClientConn{}
}

func (c *MockClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) (err error) {
	switch method {
	case "/pttbbs.api.ManService/List":
	case "/pttbbs.api.ManService/Article":
	}
	return nil
}

func (c *MockClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (cs grpc.ClientStream, err error) {
	return nil, nil
}

func (c *MockClientConn) Close() {
}
