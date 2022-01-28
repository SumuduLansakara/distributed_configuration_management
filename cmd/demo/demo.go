package main

import (
	"go_client/internal/demo"
	"go_client/pkg/utils"
	"time"
)

func main() {
	utils.InitLogging()

	system := demo.System{}
	system.Start()

	go system.StartTemperatureSensor()
	go system.StartAC()
	go system.StartDisplay()

	system.InitAcController()

	time.Sleep(time.Second * 1000)
	system.Stop()
}
