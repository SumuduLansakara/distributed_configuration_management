package prototype

import (
	"go_client/pkg/communicator"
)

type System struct {
	temperatureSensor *TemperatureSensor
	humiditySensor    *HumiditySensor
	displayUnit       *DisplayUnit
	acUnit            *AirConditioner
	acController      *AirConditionerController
}

func (s *System) Init() {
	communicator.InitClient()
	// create sensors
	s.temperatureSensor = CreateTemperatureSensor("ts-1")
	s.humiditySensor = CreateHumiditySensor("hs-1")
	// create display-unit
	s.displayUnit = CreateDisplayUnit("disp-1")
	// create air-conditioner controller
	s.acController = CreateAirConditionerController("ac-ctl-1")
	// create air-conditioner
	s.acUnit = CreateAirConditioner("ac-1")
	s.acUnit.SetParam("active", "false")
}

func (s *System) Start() {
	go s.temperatureSensor.Start()
	go s.humiditySensor.Start()
	go s.displayUnit.Start()
	go s.acUnit.Start()
	s.acController.Start()
}

func (s *System) Stop() {
	s.acController.Disconnect()
	s.acUnit.Disconnect()
	s.displayUnit.Disconnect()
	s.humiditySensor.Disconnect()
	s.temperatureSensor.Disconnect()
	communicator.DestroyClient()
}
