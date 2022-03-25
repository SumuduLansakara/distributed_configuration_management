package demo

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
	doneChan   chan interface{}
}

func (s *System) Init() {
	s.doneChan = make(chan interface{})
	communicator.InitClient()
}

func (s *System) StartTemperatureSensor() {
	ts := CreateTemperatureSensor("ts-1")
	s.components = append(s.components, ts)
	ts.Start()
}

func (s *System) StartHumiditySensor() {
	hs := CreateHumiditySensor("hs-1")
	s.components = append(s.components, hs)
	hs.Start()
}

func (s *System) StartDisplayUnit() {
	disp := CreateDisplayUnit("disp-1")
	s.components = append(s.components, disp)
	time.Sleep(1 * time.Second) // FIXME: delay till sensors connect
	disp.Start()
	<-s.doneChan
}

func (s *System) StartAC() {
	ac := CreateAirConditioner("ac-1")
	s.components = append(s.components, ac)
	ac.SetParam("active", "false")
	ac.Start()
}

func (s *System) StartACController() {
	acCtl := CreateAirConditionerController("ac-ctl-1")
	s.components = append(s.components, acCtl)
	time.Sleep(1 * time.Second) // FIXME: delay till AC connect
	acCtl.Start()
	<-s.doneChan
}

func (s *System) Stop() {
	s.doneChan <- nil
	for _, comp := range s.components {
		comp.Disconnect()
	}
	communicator.DestroyClient()
}
