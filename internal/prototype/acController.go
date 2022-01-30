package prototype

import (
	"go_client/pkg/component"
	"strconv"
)

type AirConditionerController struct {
	*DemoComponent
}

func CreateAirConditionerController(name string) *AirConditionerController {
	s := create(KindController, name)
	return &AirConditionerController{s}
}

func (c *AirConditionerController) Start() {
	// demo:
	// - a component can change parameters belonging to another component

	// get temperature sensor ID and watch its updates
	var temperatureSensor *component.RemoteComponent
	compMap := c.ListComponents(KindSensor)
	for _, stub := range compMap[KindSensor] {
		rc := stub.Spawn()
		if rc.GetParam(ParamSensorType) == ValueSensorTypeTemperatureSensor {
			temperatureSensor = rc
			break
		}
	}
	// get AC-unit ID and update its control parameters
	var airConditioner *component.RemoteComponent
	compMap = c.ListComponents(KindActuator)
	for _, stub := range compMap[KindActuator] {
		rc := stub.Spawn()
		if rc.GetParam(ParamActuatorType) == ValueActuatorTypeAirConditioner {
			airConditioner = rc
			break
		}
	}

	c.WatchParameters(temperatureSensor.Id,
		func(key, val string, isModified bool) {
			temp, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic(err)
			}
			if temp > 20 {
				if airConditioner.GetParam("active") == "false" {
					c.log("signal AC on")
					airConditioner.SetParam("active", "true")
				}
			} else if temp < 15 {
				if airConditioner.GetParam("active") == "true" {
					c.log("signal AC off")
					airConditioner.SetParam("active", "false")
				}
			}
		},
		func(key string) {
		})
}
