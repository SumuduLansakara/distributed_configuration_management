package component

import (
	"context"
	"go.uber.org/zap"
	"go_client/internal/communicator"
	"time"
)

type RemoteComponent struct {
	Component
}

func NewRemoteComponent(kind, name, id string) *RemoteComponent {
	return &RemoteComponent{
		Component: Component{Kind: kind, Name: name, Id: id},
	}
}

func (c *RemoteComponent) SetParam(key, val string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if res, err := communicator.Client.Put(ctx, path(PrefixParameters, c.Id, key), val); err != nil {
		zap.L().Panic("failed to set parameter", zap.Error(err), zap.Any("response", res))
	}
}

func (c *RemoteComponent) GetParam(key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := communicator.Client.Get(ctx, path(PrefixParameters, c.Id, key))
	cancel()
	if err != nil {
		zap.L().Panic("failed to set parameter", zap.Error(err), zap.Any("response", res))
	}
	var vals []string
	for _, v := range res.Kvs {
		vals = append(vals, string(v.Value))
	}
	if len(vals) == 0 {
		zap.L().Panic("Parameter not set", zap.String("key", key))
	}
	if len(vals) > 1 {
		zap.L().Panic("Parameter has multiple values", zap.String("key", key), zap.Any("values", vals))
	}
	return vals[0]
}

func (c *RemoteComponent) DeleteParam(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := communicator.Client.Delete(ctx, path(PrefixParameters, c.Id, key))
	cancel()
	if err != nil {
		zap.L().Panic("failed to set parameter", zap.Error(err), zap.Any("response", res))
	}
}
