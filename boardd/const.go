package boardd

import (
	"sync"

	grpc "google.golang.org/grpc"
)

var (
	conn *grpc.ClientConn
	Cli  BoardServiceClient

	cliLock sync.Mutex
)
