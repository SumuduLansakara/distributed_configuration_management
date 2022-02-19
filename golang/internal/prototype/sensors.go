package prototype

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

func initRand() {
	// Use a random seed based on the hostname to have a unique seed for each container
	hostname, err := os.Hostname()
	if err != nil {
		zap.L().Panic("failed getting hostname", zap.Error(err))
	}
	h := fnv.New64a()
	h.Write([]byte(hostname))
	hash := h.Sum64()
	seed := time.Now().UnixNano() * int64(hash)
	zap.L().Info("random seed", zap.Int64("seed", seed))
	rand.Seed(seed)
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

	initRand()
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

	initRand()
	for {
		humidity := GetHumidity()
		c.metric.Set(humidity)
		c.SetParam(ParamHumidity, fmt.Sprintf("%f", humidity))
		time.Sleep(time.Second * 2)
	}
}
