package demo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

func QueryEnvironment(key string) (float64, error) {
	resp, err := http.Get("http://environment:3100/get?key=" + key)
	if err != nil {
		zap.L().Error("failed querying environment", zap.Error(err))
		return -1, err
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("failed reading response", zap.Error(err))
		return -1, err
	}
	respMap := map[string]string{}
	jsonErr := json.Unmarshal(resBody, &respMap)
	if jsonErr != nil {
		zap.L().Error("unmarshal failed", zap.Error(err), zap.String("body", string(resBody)))
		return -1, err
	}
	v, err := strconv.ParseFloat(respMap["value"], 64)
	if err != nil {
		zap.L().Error("invalid value", zap.Error(err), zap.String("value", respMap["value"]))
		return -1, err
	}
	return v, nil
}

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

	// wait for environment to be ready
	for {
		_, err := QueryEnvironment(EnvPropertyTemperature)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 2)
	}

	for {
		time.Sleep(time.Second * 2)
		temperature, err := QueryEnvironment(EnvPropertyTemperature)
		if err != nil {
			continue
		}
		c.metric.Set(temperature)
		c.SetParam(ParamTemperature, fmt.Sprintf("%f", temperature))
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

	// wait for environment to be ready
	for {
		_, err := QueryEnvironment(EnvPropertyHumidity)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 2)
	}

	for {
		time.Sleep(time.Second * 2)
		humidity, err := QueryEnvironment(EnvPropertyHumidity)
		if err != nil {
			continue
		}
		c.metric.Set(humidity)
		c.SetParam(ParamHumidity, fmt.Sprintf("%f", humidity))
	}
}
