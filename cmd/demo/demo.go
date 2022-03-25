package main

import (
	"go_client/internal/demo"
	"go_client/pkg/utils"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func startMetricServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":3001", nil)
}

func main() {
	utils.InitLogging()

	go startMetricServer()

	logger := zap.L()
	logger.Debug("component started", zap.Any("args", os.Args))

	if len(os.Args) != 2 {
		zap.L().Panic("Invalid argument count")
	}

	system := demo.System{}
	zap.L().Debug("starting", zap.String("name", os.Args[1]))

	// following are blocking calls
	switch os.Args[1] {
	case "environment":
		demo.StartEnvironment()
	case "temperature-sensor":
		system.Init()
		system.StartTemperatureSensor()
	case "humidity-sensor":
		system.Init()
		system.StartHumiditySensor()
	case "display":
		system.Init()
		system.StartDisplayUnit()
	case "airconditioner":
		system.Init()
		system.StartAC()
	case "airconditioner-controller":
		system.Init()
		system.StartACController()
	default:
		zap.L().Panic("Unsupported component name", zap.String("name", os.Args[1]))
	}
	zap.L().Debug("stopping", zap.String("name", os.Args[1]))
}
