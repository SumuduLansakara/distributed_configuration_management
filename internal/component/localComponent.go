package component

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"go_client/internal/communicator"
	"strings"
	"time"
)

type LocalComponent struct {
	Component
	configs   map[string]string
	connected bool
}

func (c *LocalComponent) Test() {
	zap.L().Info("test", zap.Any("configs", c.configs))
}

func NewLocalComponent(kind, name string) (*LocalComponent, error) {
	newUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "failed creating UUID")
	}
	return &LocalComponent{
		Component: Component{Kind: kind, Name: name, InstanceID: newUuid.String()},
		configs:   map[string]string{},
	}, nil
}

func (c *LocalComponent) paramKey(path []string) string {
	return strings.Join(path, "/")
}

func (c *LocalComponent) metadataKey(path []string) string {
	return strings.Join(path, "/")
}

func (c *LocalComponent) componentKey(path []string) string {
	return strings.Join(path, "/")
}

func (c *LocalComponent) Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ops := []clientV3.Op{
		clientV3.OpPut(c.metadataKey([]string{PrefixMetadata, c.InstanceID, "name"}), c.Name),
		clientV3.OpPut(c.metadataKey([]string{PrefixMetadata, c.InstanceID, "kind"}), c.Kind),
		clientV3.OpPut(c.componentKey([]string{PrefixComponents, c.Kind, c.InstanceID}), ""),
	}
	for _, op := range ops {
		if _, err := communicator.Client.Do(ctx, op); err != nil {
			zap.L().Panic("connection failed", zap.Error(err))
		}
	}
	c.connected = true
}

func (c *LocalComponent) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ops := []clientV3.Op{
		clientV3.OpDelete(c.metadataKey([]string{PrefixMetadata, c.InstanceID}), clientV3.WithPrefix()),
		clientV3.OpDelete(c.componentKey([]string{PrefixComponents, c.Kind, c.InstanceID}), clientV3.WithPrefix()),
	}
	for _, op := range ops {
		if _, err := communicator.Client.Do(ctx, op); err != nil {
			zap.L().Panic("connection failed", zap.Error(err))
		}
	}
	c.connected = true
}

func (c *LocalComponent) SetParam(key, val string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if res, err := communicator.Client.Put(ctx, c.paramKey([]string{PrefixParameters, c.InstanceID, key}), val); err != nil {
		zap.L().Panic("failed to set parameter", zap.Error(err), zap.Any("response", res))
	}
	c.configs[key] = val
}

func (c *LocalComponent) GetParam(key string) string {
	val, ok := c.configs[key]
	if !ok {
		zap.L().Panic("Parameter not set", zap.String("key", key))
	}
	return val
}

func (c *LocalComponent) DeleteParam(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := communicator.Client.Delete(ctx, c.paramKey([]string{PrefixParameters, c.InstanceID, key}))
	cancel()
	if err != nil {
		zap.L().Panic("failed to set parameter", zap.Error(err), zap.Any("response", res))
	}
	delete(c.configs, key)
}

func (c *LocalComponent) ReloadParam(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := communicator.Client.Get(ctx, c.paramKey([]string{PrefixParameters, c.InstanceID, key}))
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
	c.configs[key] = vals[0]
}

func (c *LocalComponent) ReloadAllParams() {
	prefix := strings.Join([]string{PrefixParameters, c.InstanceID}, "/")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := communicator.Client.Get(ctx, prefix, clientV3.WithPrefix())
	cancel()
	if err != nil {
		zap.L().Panic("failed reloading all params", zap.Error(err), zap.Any("response", res))
	}
	zap.L().Info("all params", zap.Any("response", res))
}

func (c *LocalComponent) ListComponents(kind string) map[string][]string {
	path := strings.Join([]string{PrefixComponents, kind}, "/")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := communicator.Client.Get(ctx, path, clientV3.WithPrefix())
	cancel()
	if err != nil {
		zap.L().Panic("failed listing components", zap.Error(err), zap.Any("response", res))
	}
	components := map[string][]string{}
	for _, v := range res.Kvs {
		tokens := strings.Split(string(v.Key), "/")
		thisKind := tokens[2]
		thisUuid := tokens[3]
		if _, ok := components[thisKind]; !ok {
			components[thisKind] = []string{}
		}
		components[thisKind] = append(components[thisKind], thisUuid)
	}
	zap.L().Info("components list", zap.Any("components", components))
	return components
}
