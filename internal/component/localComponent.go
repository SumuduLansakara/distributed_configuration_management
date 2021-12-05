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

func path(path ...string) string {
	return strings.Join(path, "/")
}

func (c *LocalComponent) PrettyName() string {
	return strings.Join([]string{c.Kind, c.Name, c.InstanceID}, ":")
}

func (c *LocalComponent) Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ops := []clientV3.Op{
		clientV3.OpPut(path(PrefixMetadata, c.InstanceID, "name"), c.Name),
		clientV3.OpPut(path(PrefixMetadata, c.InstanceID, "kind"), c.Kind),
		clientV3.OpPut(path(PrefixComponents, c.Kind, c.InstanceID), ""),
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
		clientV3.OpDelete(path(PrefixMetadata, c.InstanceID), clientV3.WithPrefix()),
		clientV3.OpDelete(path(PrefixComponents, c.Kind, c.InstanceID), clientV3.WithPrefix()),
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
	if res, err := communicator.Client.Put(ctx, path(PrefixParameters, c.InstanceID, key), val); err != nil {
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
	res, err := communicator.Client.Delete(ctx, path(PrefixParameters, c.InstanceID, key))
	cancel()
	if err != nil {
		zap.L().Panic("failed to set parameter", zap.Error(err), zap.Any("response", res))
	}
	delete(c.configs, key)
}

func (c *LocalComponent) ReloadParam(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := communicator.Client.Get(ctx, path(PrefixParameters, c.InstanceID, key))
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := communicator.Client.Get(ctx, path(PrefixComponents, kind), clientV3.WithPrefix())
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
			components[thisKind] = []string{thisUuid}
		} else {
			components[thisKind] = append(components[thisKind], thisUuid)
		}
	}
	zap.L().Info("components list", zap.Any("components", components))
	return components
}

func (c *LocalComponent) WatchComponents(kind string) {
	go func() {
		watchChan := communicator.Client.Watch(context.Background(), path(PrefixComponents, kind), clientV3.WithPrefix())
		zap.L().Debug("watch begin", zap.String("watcher", c.PrettyName()), zap.String("target kind", kind))
		for wresp := range watchChan {
			for _, ev := range wresp.Events {
				tokens := strings.Split(string(ev.Kv.Key), "/")
				thisKind := tokens[2]
				thisUuid := tokens[3]
				switch ev.Type {
				case clientV3.EventTypePut:
					zap.L().Info("component created", zap.String("kind", thisKind), zap.String("uuid", thisUuid))
				case clientV3.EventTypeDelete:
					zap.L().Info("component deleted", zap.String("kind", thisKind), zap.String("uuid", thisUuid))
				default:
					zap.L().Debug("unknown component action", zap.Any("type", ev.Type), zap.String("kind", thisKind), zap.String("uuid", thisUuid))
				}
			}
		}
		zap.L().Debug("watch end", zap.String("watcher", c.PrettyName()), zap.String("target kind", kind))
	}()
}
