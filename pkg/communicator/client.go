package communicator

import (
	"context"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

var Client *clientV3.Client

func InitClient() {
	var err error
	Client, err = clientV3.New(clientV3.Config{
		Endpoints:   []string{"127.0.0.1:12379", "127.0.0.1:12380"},
		DialTimeout: 5 * time.Second,
	})
	checkError(err)
}

func DestroyClient() {
	Client.Close()
}

func checkError(err error) {
	if err != nil {
		logger := zap.L()
		switch err {
		case context.Canceled:
			logger.Fatal("ctx is canceled by another routine", zap.Error(err))
		case context.DeadlineExceeded:
			logger.Fatal("ctx is attached with a deadline is exceeded", zap.Error(err))
		case rpctypes.ErrEmptyKey:
			logger.Fatal("client-side error", zap.Error(err))
		default:
			logger.Fatal("bad cluster endpoints, which are not etcd servers", zap.Error(err))
		}
	}
}
