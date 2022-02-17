package communicator

import (
	"context"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var Client *clientV3.Client

func InitClient() {
	var connectionErr error
	attempts := 0
	for ; attempts < 10; attempts++ {
		Client, connectionErr = clientV3.New(clientV3.Config{
			Endpoints:   []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"},
			DialTimeout: 5 * time.Second,
		})
		if connectionErr != nil {
			zap.L().Debug("client creation failed", zap.Int("retryCount", attempts))
			continue
		}
		//
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, connectionErr = Client.Get(ctx, "/")
		if connectionErr != nil {
			zap.L().Debug("test query failed", zap.Int("retryCount", attempts))
			continue
		}
		break
	}
	if connectionErr == nil {
		zap.L().Debug("successfully connected to etcd cluster", zap.Int("retryCount", attempts))
	} else {
		zap.L().Panic("failed creating etcd client", zap.Error(connectionErr), zap.Int("retryCount", attempts))
	}
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
