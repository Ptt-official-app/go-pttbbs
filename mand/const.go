package mand

import (
	sync "sync"

	grpc "google.golang.org/grpc"
)

var (
	conn *grpc.ClientConn
	Cli  ManServiceClient

	cliLock sync.Mutex
)
