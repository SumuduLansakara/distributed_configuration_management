package prototype

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// TemperatureSensor mocks a temperature sensor by periodically updating random temperature readings
type TemperatureSensor struct {
	*DemoComponent
	metric prometheus.Gauge
}

func CreateTemperatureSensor(name string) *TemperatureSensor {
	c := create(KindSensor, name)
	c.SetParam(ParamSensorType, ValueSensorTypeTemperatureSensor)
	c.SetParam(ParamTemperature, "0") // initial value
	return &TemperatureSensor{
		DemoComponent: c,
		metric: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "current_temparature",
			Help: "Current environment temperature",
		}),
	}
}

func (c *TemperatureSensor) Start() {
	// demo:
	// - a component can change its own parameters

	for {
		temperature := GetTemperature()
		c.metric.Set(temperature)
		c.SetParam(ParamTemperature, fmt.Sprintf("%f", temperature))
		time.Sleep(time.Second * 2)
	}
}

// HumiditySensor mocks a humidity sensor by periodically updating random humidity readings
type HumiditySensor struct {
	*DemoComponent
	metric prometheus.Gauge
}

func CreateHumiditySensor(name string) *HumiditySensor {
	c := create(KindSensor, name)
	c.SetParam(ParamSensorType, ValueSensorTypeHumiditySensor)
	c.SetParam(ParamHumidity, "0") // initial value
	return &HumiditySensor{
		DemoComponent: c,
		metric: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "current_humidity",
			Help: "Current environment humidity",
		}),
	}
}

func (c *HumiditySensor) Start() {
	// demo:
	// - a component can change its own parameters

	for {
		humidity := GetHumidity()
		c.metric.Set(humidity)
		c.SetParam(ParamHumidity, fmt.Sprintf("%f", humidity))
		time.Sleep(time.Second * 2)
	}
}
