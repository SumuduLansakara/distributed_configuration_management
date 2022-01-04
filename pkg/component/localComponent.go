package component

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"go_client/pkg/communicator"
	"strings"
	"sync"
	"time"
)

type LocalComponent struct {
	Component
	configs         map[string]string
	watchCancellers map[string]context.CancelFunc
	connected       bool
	configLock      sync.RWMutex
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
		Component:       Component{Kind: kind, Name: name, Id: newUuid.String()},
		configs:         map[string]string{},
		watchCancellers: map[string]context.CancelFunc{},
		connected:       false,
		configLock:      sync.RWMutex{},
	}, nil
}

func path(path ...string) string {
	return strings.Join(path, "/")
}

func (c *LocalComponent) PrettyName() string {
	return strings.Join([]string{c.Kind, c.Name, c.Id}, ":")
}

func (c *LocalComponent) Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ops := []clientV3.Op{
		clientV3.OpPut(path(PrefixMetadata, c.Id, "name"), c.Name),
		clientV3.OpPut(path(PrefixMetadata, c.Id, "kind"), c.Kind),
		clientV3.OpPut(path(PrefixComponents, c.Kind, c.Id), ""),
	}
	for _, op := range ops {
		if _, err := communicator.Client.Do(ctx, op); err != nil {
			zap.L().Panic("connection failed", zap.Error(err))
		}
	}
	c.connected = true
	// listen to changes of my params
	c.WatchParameters(c.Id,
		func(key, val string, isModified bool) {
			c.configLock.Lock()
			defer c.configLock.Unlock()
			c.configs[key] = val
		},
		func(key string) {
			c.configLock.Lock()
			defer c.configLock.Unlock()
			delete(c.configs, key)
		},
	)
}

func (c *LocalComponent) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// remove watchers
	for _, cancelWatch := range c.watchCancellers {
		cancelWatch()
	}
	c.watchCancellers = map[string]context.CancelFunc{}
	// remove configs
	ops := []clientV3.Op{
		clientV3.OpDelete(path(PrefixComponents, c.Kind, c.Id), clientV3.WithPrefix()),
		clientV3.OpDelete(path(PrefixMetadata, c.Id), clientV3.WithPrefix()),
		clientV3.OpDelete(path(PrefixParameters, c.Id), clientV3.WithPrefix()),
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
	if res, err := communicator.Client.Put(ctx, path(PrefixParameters, c.Id, key), val); err != nil {
		zap.L().Panic("failed to set parameter", zap.Error(err), zap.Any("response", res))
	}
	c.configLock.Lock()
	defer c.configLock.Unlock()
}

func (c *LocalComponent) IsParamSet(key string) bool {
	c.configLock.RLock()
	defer c.configLock.RUnlock()
	_, ok := c.configs[key]
	return ok
}

func (c *LocalComponent) GetParam(key string) string {
	var val string
	var ok bool
	func() {
		c.configLock.RLock()
		defer c.configLock.RUnlock()
		val, ok = c.configs[key]
	}()
	if !ok {
		zap.L().Panic("Parameter not set", zap.String("key", key))
	}
	return val
}

func (c *LocalComponent) DeleteParam(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := communicator.Client.Delete(ctx, path(PrefixParameters, c.Id, key))
	cancel()
	if err != nil {
		zap.L().Panic("failed to set parameter", zap.Error(err), zap.Any("response", res))
	}
	c.configLock.Lock()
	defer c.configLock.Unlock()
}

func (c *LocalComponent) ReloadParam(key string) {
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
	c.configLock.Lock()
	defer c.configLock.Unlock()
	c.configs[key] = vals[0]
}

func (c *LocalComponent) ReloadAllParams() {
	prefix := strings.Join([]string{PrefixParameters, c.Id}, "/")
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

func (c *LocalComponent) WatchComponents(kind string, onConnected, onDeleted func(stub *RemoteComponentStub)) {
	if !c.connected {
		zap.L().Warn("disconnected component can't add a component watch")
		return
	}
	watchPath := path(PrefixComponents, kind)
	if _, ok := c.watchCancellers[watchPath]; ok {
		zap.L().Warn("ignoring duplicated component watch", zap.String("watchPath", watchPath))
		return
	}
	watchCtx, cancelWatch := context.WithCancel(context.Background())
	watchChan := communicator.Client.Watch(watchCtx, watchPath, clientV3.WithPrefix())
	c.watchCancellers[watchPath] = cancelWatch
	go func() {
		for wresp := range watchChan {
			for _, ev := range wresp.Events {
				tokens := strings.Split(string(ev.Kv.Key), "/")
				thisKind := tokens[2]
				thisId := tokens[3]
				switch ev.Type {
				case clientV3.EventTypePut:
					if onConnected != nil {
						onConnected(NewRemoteComponentStub(thisKind, thisId))
					}
				case clientV3.EventTypeDelete:
					if onDeleted != nil {
						onDeleted(NewRemoteComponentStub(thisKind, thisId))
					}
				default:
					zap.L().Warn("unknown component event", zap.Any("type", ev.Type), zap.String("kind", thisKind), zap.String("id", thisId))
				}
			}
		}
	}()
}

func (c *LocalComponent) WatchParameters(componentID string, onSet func(key, val string, isModified bool), onDeleted func(key string)) {
	if !c.connected {
		zap.L().Warn("disconnected component can't add a parameter watch")
		return
	}
	watchPath := path(PrefixParameters, componentID)
	if _, ok := c.watchCancellers[watchPath]; ok {
		zap.L().Warn("ignoring duplicated parameter watch", zap.String("watchPath", watchPath))
		return
	}
	watchCtx, cancelWatch := context.WithCancel(context.Background())
	watchChan := communicator.Client.Watch(watchCtx, watchPath, clientV3.WithPrefix())
	c.watchCancellers[watchPath] = cancelWatch
	go func() {
		for watchResponse := range watchChan {
			for _, ev := range watchResponse.Events {
				tokens := strings.Split(string(ev.Kv.Key), "/")
				key := tokens[3]
				val := string(ev.Kv.Value)
				switch ev.Type {
				case clientV3.EventTypePut:
					if onSet != nil {
						onSet(key, val, ev.IsModify())
					}
				case clientV3.EventTypeDelete:
					if onDeleted != nil {
						onDeleted(key)
					}
				default:
					zap.L().Warn("unknown param event", zap.Any("type", ev.Type), zap.String("key", key), zap.String("val", val))
				}
			}
		}
	}()
}
