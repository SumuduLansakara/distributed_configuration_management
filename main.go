package main

import (
	"go.uber.org/zap"
	"go_client/internal/communicator"
	"go_client/internal/component"
	"go_client/internal/utils"
	"time"
)

func init() {
	utils.InitLogging()
}

func main() {
	communicator.InitClient()
	defer communicator.DestroyClient()

	comp, err := component.NewLocalComponent("mykind", "myname")
	if err != nil {
		zap.L().Panic("construction failed", zap.Error(err))
	}
	comp.Connect()
	zap.L().Info("components list", zap.Any("components", comp.ListComponents("")))

	comp.WatchComponents(comp.Kind, func(stub *component.RemoteComponentStub) {
		comp.WatchParameters(stub, func(key, val string) {
			zap.L().Info("param set", zap.String("key", key), zap.String("val", val))
		}, nil)
	}, nil)
	//comp.ReloadAllParams()
	//comp.Test()

	comp2, _ := component.NewLocalComponent("mykind", "myname2")
	comp2.Connect()
	comp2.SetParam("foo", "bar")
	time.Sleep(time.Second * 3)
	comp2.Disconnect()
	comp.Disconnect()

	zap.L().Info("components list", zap.Any("components", comp.ListComponents("")))
}
