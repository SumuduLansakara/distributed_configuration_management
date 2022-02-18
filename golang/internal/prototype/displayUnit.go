package prototype

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type DisplayUnit struct {
	*DemoComponent
	metric prometheus.Gauge
}

func CreateDisplayUnit(name string) *DisplayUnit {
	s := create(KindDisplay, name)
	return &DisplayUnit{
		DemoComponent: s,
		metric: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "display_value",
			Help: "Current display value",
		}),
	}
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
					valF, err := strconv.ParseFloat(val, 64)
					if err != nil {
						zap.L().Panic("failed parsing sensor reading", zap.Error(err))
					}
					c.log(fmt.Sprintf("temperature: %f", valF))
					c.metric.Set(valF)
				},
				func(key string) {
				})
		}
	}
}
