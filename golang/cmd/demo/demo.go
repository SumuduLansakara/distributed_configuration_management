package main

import (
	"go_client/internal/prototype"
	"go_client/pkg/utils"
	"net/http"
	"os"
	"time"

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

	// stop the simulation after 10 minutes
	time.Sleep(time.Second * 600)
	system.Stop()
}
