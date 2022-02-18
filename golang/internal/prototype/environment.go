package prototype

import (
	"math/rand"
)

var (
	temperature = 15.0
	humidity    = 30.0
)

func GetTemperature() float64 {
	temperature += rand.Float64()*10 - 4
	return temperature
}

func GetHumidity() float64 {
	humidity += rand.Float64()*10 - 4
	return humidity
}
