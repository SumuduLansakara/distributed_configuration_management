package demo

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type AirConditioner struct {
	*Component
	isActive bool
	metric   prometheus.Gauge
}

func CreateAirConditioner(name string) *AirConditioner {
	c := create(KindActuator, name)
	c.SetParam(ParamActuatorType, ValueActuatorTypeAirConditioner)
	c.SetParam(ParamACState, ValueACStateInactive)
	return &AirConditioner{
		Component: c,
		isActive:  false,
		metric: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "airconditioner_state",
			Help: "Current state of the air-conditioner",
		}),
	}
}

func (c *AirConditioner) Start() {
	// demo:
	// - a component can see changes to local parameters by external components
	for {
		time.Sleep(time.Second * 1) // check if my parameters are changed (by AC controller)
		if !c.isActive && c.GetParam(ParamACState) == ValueACStateActive {
			c.isActive = true
			c.log("turned-on")
			c.metric.Set(1)
		} else if c.isActive && c.GetParam(ParamACState) == ValueACStateInactive {
			c.isActive = false
			c.log("turned-off")
			c.metric.Set(0)
		}
	}
}
