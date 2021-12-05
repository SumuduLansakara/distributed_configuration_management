package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed initializing logger")
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)
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

func main() {
	cli, err := clientV3.New(clientV3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	checkError(err)
	defer cli.Close()

	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		putResp, err := cli.Put(ctx, "testkey1", "testval1")
		cancel()
		checkError(err)
		zap.L().Info("put success", zap.Any("resp", putResp))
	}

	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		getResp, err := cli.Get(ctx, "testkey1")
		cancel()
		checkError(err)
		zap.L().Info("get success", zap.Any("resp", getResp))
		for i, ev := range getResp.Kvs {
			zap.L().Info("KV", zap.Int("i", i), zap.ByteString("key", ev.Key), zap.ByteString("val", ev.Value))
		}
	}
}
