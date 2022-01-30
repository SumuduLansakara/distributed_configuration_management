package main

import (
	"go_client/internal/prototype"
	"go_client/pkg/utils"
	"time"
)

func main() {
	utils.InitLogging()

	system := prototype.System{}
	system.Init()
	system.Start()

	// stop the simulation after 1 minute
	time.Sleep(time.Second * 60)
	system.Stop()
}
