package component

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"go_client/internal/communicator"
	"strings"
	"time"
)

type RemoteComponentStub struct {
	Kind string
	Id   string
}

func NewRemoteComponentStub(kind, id string) *RemoteComponentStub {
	return &RemoteComponentStub{Kind: kind, Id: id}
}

func (c *RemoteComponentStub) Spawn() *RemoteComponent {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := communicator.Client.Get(ctx, path(PrefixMetadata, c.Id), clientv3.WithPrefix())
	cancel()
	if err != nil {
		zap.L().Panic("failed reading meta data", zap.Error(err))
	}
	for _, v := range res.Kvs {
		tokens := strings.Split(string(v.Key), "/")
		key := tokens[3]
		if key == "name" {
			return NewRemoteComponent(c.Kind, string(v.Value), c.Id)
		}
	}
	zap.L().Panic("remote component not found")
	return nil
}
