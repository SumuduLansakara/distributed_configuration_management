package main

import (
	"go_client/internal/prototype"
	"go_client/pkg/utils"
	"os"
	"time"

	"go.uber.org/zap"
)

func main() {
	utils.InitLogging()
	logger := zap.L()
	logger.Debug("component started", zap.Any("args", os.Args))

	if len(os.Args) != 2 {
		zap.L().Panic("Invalid argument count")
	}

	system := prototype.System{}
	system.Init()

	switch os.Args[1] {
	case "temperature-sensor":
		system.StartTemperatureSensor()
	case "humidity-sensor":
		system.StartHumiditySensor()
	case "display":
		system.StartDisplayUnit()
	case "airconditioner":
		system.StartAC()
	case "airconditioner-controller":
		system.StartACController()
	default:
		zap.L().Panic("Unsupported component name", zap.String("name", os.Args[1]))
	}
	zap.L().Debug("started", zap.String("name", os.Args[1]))

	// stop the simulation after 1 minute
	time.Sleep(time.Second * 60)
	system.Stop()
}
