package boardd

import (
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init(isLocked bool) (err error) {
	if !isLocked {
		cliLock.Lock()
		defer cliLock.Unlock()
	}

	if cli != nil {
		return ErrCliAlreadyInit
	}

	if IsTest {
		mockConn := NewMockClientConn()
		cli = NewBoardServiceClient(mockConn)
	} else {
		conn, err = grpc.Dial("localhost:5150", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}

		cli = NewBoardServiceClient(conn)
	}

	return nil
}

func Finalize(isLocked bool) {
	if !isLocked {
		cliLock.Lock()
		defer cliLock.Unlock()
	}

	defer func() {
		if conn != nil {
			conn.Close()
			conn = nil
		}
	}()

	defer func() {
		cli = nil
	}()
}

func Reset() (err error) {
	cliLock.Lock()
	defer cliLock.Unlock()

	Finalize(true)

	return Init(true)
}
