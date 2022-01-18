package boardd

import (
	"sync"

	grpc "google.golang.org/grpc"
)

var (
	conn *grpc.ClientConn
	cli  BoardServiceClient

	cliLock sync.Mutex
)
