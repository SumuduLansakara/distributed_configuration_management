package prototype

import (
	"time"
)

type AirConditioner struct {
	*DemoComponent
	isActive bool
}

func CreateAirConditioner(name string) *AirConditioner {
	c := create(KindActuator, name)
	c.SetParam(ParamActuatorType, ValueActuatorTypeAirConditioner)
	c.SetParam(ParamACState, ValueACStateInactive)
	return &AirConditioner{DemoComponent: c, isActive: false}
}

func (c *AirConditioner) Start() {
	// demo:
	// - a component can see changes to local parameters by external components
	for {
		time.Sleep(time.Second * 1) // check if my parameters are changed (by AC controller)
		if !c.isActive && c.GetParam(ParamACState) == ValueACStateActive {
			c.isActive = true
			c.log("turned-on")
		} else if c.isActive && c.GetParam(ParamACState) == ValueACStateInactive {
			c.isActive = false
			c.log("turned-off")
		}
	}
}
