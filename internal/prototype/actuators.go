package prototype

import (
	"time"
)

type AirConditioner struct {
	*DemoComponent
}

func CreateAirConditioner(name string) *AirConditioner {
	c := create(KindActuator, name)
	c.SetParam(ParamActuatorType, ValueActuatorTypeAirConditioner)
	return &AirConditioner{c}
}

func (c *AirConditioner) Start() {
	// demo:
	// - a component can see changes to local parameters by external components
	isActive := false
	for {
		time.Sleep(time.Second * 1) // check if my parameters are changed (by AC controller)
		if !isActive && c.GetParam("active") == "true" {
			isActive = true
			c.log("turned-on")
		} else if isActive && c.GetParam("active") == "false" {
			isActive = false
			c.log("turned-off")
		}
	}
}
