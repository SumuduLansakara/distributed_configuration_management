package prototype

import (
	"fmt"
)

type DisplayUnit struct {
	*DemoComponent
}

func CreateDisplayUnit(name string) *DisplayUnit {
	s := create(KindDisplay, name)
	return &DisplayUnit{s}
}

func (c *DisplayUnit) Start() {
	// demo:
	// - a component can list all the components of selected type in the system
	// - a component can watch parameter changes of interested components
	compMap := c.ListComponents(KindSensor)
	for _, stub := range compMap[KindSensor] {
		rc := stub.Spawn()
		if rc.GetParam(ParamSensorType) == ValueSensorTypeTemperatureSensor {
			c.WatchParameters(rc.Id,
				func(key, val string, isModified bool) {
					c.log(fmt.Sprintf("temperature: %s", val))
				},
				func(key string) {
				})
		}
	}
}
