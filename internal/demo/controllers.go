package demo

import (
	"go_client/pkg/component"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type AirConditionerController struct {
	*DemoComponent
	metric prometheus.Gauge
}

func CreateAirConditionerController(name string) *AirConditionerController {
	s := create(KindController, name)
	return &AirConditionerController{
		DemoComponent: s,
		metric: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "airconditioner_controller_state",
			Help: "Current state of the air-conditioner controller",
		}),
	}
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
			if temp > 25 {
				if airConditioner.GetParam(ParamACState) == ValueACStateInactive {
					c.log("signal AC on")
					c.metric.Set(1)
					airConditioner.SetParam(ParamACState, ValueACStateActive)
				}
			} else if temp < 20 {
				if airConditioner.GetParam(ParamACState) == ValueACStateActive {
					c.log("signal AC off")
					c.metric.Set(0)
					airConditioner.SetParam(ParamACState, ValueACStateInactive)
				}
			}
		},
		func(key string) {
		})
}
