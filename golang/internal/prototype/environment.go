package prototype

import (
	"math/rand"
)

var (
	temperature = 15.0
	humidity    = 30.0
)

func getRandomInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GetTemperature() float64 {
	temperature += getRandomInRange(-5, 5)
	return temperature
}

func GetHumidity() float64 {
	humidity += getRandomInRange(-1, 1)
	if humidity < 0 {
		humidity = 0
	} else if humidity > 100 {
		humidity = 100
	}
	return humidity
}
