package main

import (
	"go.uber.org/zap"
	"go_client/internal/communicator"
	"go_client/internal/component"
	"go_client/internal/utils"
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
	componentList := comp.ListComponents("")
	zap.L().Info("components list", zap.Any("components", componentList))
	//comp.ReloadAllParams()
	//comp.Test()
	// comp.Disconnect()
}
