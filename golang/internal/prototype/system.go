package prototype

import (
	"go_client/pkg/communicator"
	"time"
)

type ComponentI interface {
	Start()
	Disconnect()
}

type System struct {
	components []ComponentI
}

func (s *System) Init() {
	communicator.InitClient()
}

func (s *System) StartTemperatureSensor() {
	ts := CreateTemperatureSensor("ts-1")
	s.components = append(s.components, ts)
	go ts.Start()
}

func (s *System) StartHumiditySensor() {
	hs := CreateHumiditySensor("hs-1")
	s.components = append(s.components, hs)
	go hs.Start()
}

func (s *System) StartDisplayUnit() {
	disp := CreateDisplayUnit("disp-1")
	s.components = append(s.components, disp)
	time.Sleep(1 * time.Second) // FIXME: delay till sensors connect
	go disp.Start()
}

func (s *System) StartAC() {
	ac := CreateAirConditioner("ac-1")
	s.components = append(s.components, ac)
	ac.SetParam("active", "false")
	go ac.Start()
}

func (s *System) StartACController() {
	acCtl := CreateAirConditionerController("ac-ctl-1")
	s.components = append(s.components, acCtl)
	time.Sleep(1 * time.Second) // FIXME: delay till AC connect
	acCtl.Start()
}

func (s *System) Stop() {
	for _, comp := range s.components {
		comp.Disconnect()
	}
	communicator.DestroyClient()
}
