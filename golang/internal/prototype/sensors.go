package prototype

import (
	"time"
)

// TemperatureSensor mocks a temperature sensor by periodically updating random temperature readings
type TemperatureSensor struct {
	*DemoComponent
}

func CreateTemperatureSensor(name string) *TemperatureSensor {
	c := create(KindSensor, name)
	c.SetParam(ParamSensorType, ValueSensorTypeTemperatureSensor)
	c.SetParam(ParamTemperature, "0") // initial value
	return &TemperatureSensor{c}
}

func (c *TemperatureSensor) Start() {
	// demo:
	// - a component can change its own parameters
	for {
		c.SetParam(ParamTemperature, GetTemperature())
		time.Sleep(time.Second * 2)
	}
}

// HumiditySensor mocks a humidity sensor by periodically updating random humidity readings
type HumiditySensor struct {
	*DemoComponent
}

func CreateHumiditySensor(name string) *HumiditySensor {
	c := create(KindSensor, name)
	c.SetParam(ParamSensorType, ValueSensorTypeHumiditySensor)
	c.SetParam(ParamHumidity, "0") // initial value
	return &HumiditySensor{c}
}

func (c *HumiditySensor) Start() {
	// demo:
	// - a component can change its own parameters
	for {
		c.SetParam(ParamHumidity, GetHumidity())
		time.Sleep(time.Second * 2)
	}
}
