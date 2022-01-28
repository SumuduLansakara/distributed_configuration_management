package demo

import (
	"fmt"
	"go_client/pkg/communicator"
	"go_client/pkg/component"
	"math/rand"
	"strconv"
	"time"
)

type Component struct {
	comp *component.LocalComponent
}

func create(kind, name string) *Component {
	comp, err := component.NewLocalComponent(kind, name)
	if err != nil {
		panic(err)
	}
	comp.Connect()
	return &Component{comp: comp}
}

func (d *Component) log(msg string) {
	fmt.Printf("[%s] %s\n", d.comp.Name, msg)
}

type System struct {
	humiditySensor    *Component
	temperatureSensor *Component
	displayUnit       *Component
	acController      *Component
	acUnit            *Component
}

func (s *System) Start() {
	communicator.InitClient()
	// create sensors
	s.humiditySensor = create("humidity-sensor", "hs-1")
	s.temperatureSensor = create("temperature-sensor", "ts-1")
	s.temperatureSensor.comp.SetParam("temperature", "10.0")
	// create display-unit
	s.displayUnit = create("display", "display-1")
	// create air-conditioner controller
	s.acController = create("controller", "ac-controller-1")
	// create air-conditioner
	s.acUnit = create("actuator", "ac-1")
	s.acUnit.comp.SetParam("active", "false")
}

func (s *System) Stop() {
	s.acUnit.comp.Disconnect()
	s.acController.comp.Disconnect()
	s.displayUnit.comp.Disconnect()
	s.temperatureSensor.comp.Disconnect()
	s.humiditySensor.comp.Disconnect()
	communicator.DestroyClient()
}

func (s *System) StartTemperatureSensor() {
	// demo:
	// - a component can change its own parameters
	temp := 10.0
	for {
		temp += rand.Float64()*10 - 4
		s.temperatureSensor.comp.SetParam("temperature", fmt.Sprintf("%f", temp))
		time.Sleep(time.Second * 2)
	}
}

func (s *System) StartDisplay() {
	// demo:
	// - a component can list all the components of selected type in the system
	// - a component can watch parameter changes of interested components
	compMap := s.displayUnit.comp.ListComponents("temperature-sensor")
	lastReading := ""
	for _, compId := range compMap["temperature-sensor"] {
		s.displayUnit.comp.WatchParameters(compId,
			func(key, val string, isModified bool) {
				if val != lastReading {
					lastReading = val
					s.displayUnit.log(fmt.Sprintf("temperature: %s", lastReading))
				}
			},
			func(key string) {
			})
	}
}

func (s *System) InitAcController() {
	// demo:
	// - a component can change parameters belonging to another component
	s.acController.comp.WatchParameters(s.temperatureSensor.comp.Id,
		func(key, val string, isModified bool) {
			temp, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic(err)
			}
			if temp > 20 {
				if s.acUnit.comp.GetParam("active") == "false" {
					s.acController.log("signal AC on")
					s.acUnit.comp.SetParam("active", "true")
				}
			} else if temp < 15 {
				if s.acUnit.comp.GetParam("active") == "true" {
					s.acController.log("signal AC off")
					s.acUnit.comp.SetParam("active", "false")
				}
			}
		},
		func(key string) {
		})
}

func (s *System) StartAC() {
	// demo:
	// - a component can see changes to local parameters by external components
	isActive := false
	for {
		time.Sleep(time.Second * 1) // check if my parameters are changed (by AC controller)
		if !isActive && s.acUnit.comp.GetParam("active") == "true" {
			isActive = true
			s.acUnit.log("turned-on")
		} else if isActive && s.acUnit.comp.GetParam("active") == "false" {
			isActive = false
			s.acUnit.log("turned-off")
		}
	}
}
