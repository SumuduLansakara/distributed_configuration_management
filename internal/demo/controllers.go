package demo

import (
	"go_client/pkg/component"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	ThresholdAcTurnOn  = 25.0
	ThresholdAcTurnOff = 20.0
)

type AirConditionerController struct {
	*Component
	metric prometheus.Gauge
}

func CreateAirConditionerController(name string) *AirConditionerController {
	s := create(KindController, name)
	return &AirConditionerController{
		Component: s,
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
			if temp > ThresholdAcTurnOn {
				if airConditioner.GetParam(ParamACState) == ValueACStateInactive {
					c.log("signalling AC to turn on")
					airConditioner.SetParam(ParamACState, ValueACStateActive)
					c.metric.Set(1)
				}
			} else if temp < ThresholdAcTurnOff {
				if airConditioner.GetParam(ParamACState) == ValueACStateActive {
					c.log("signalling AC to turn off")
					airConditioner.SetParam(ParamACState, ValueACStateInactive)
					c.metric.Set(0)
				}
			}
		},
		func(key string) {
		})
}
