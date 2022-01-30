package prototype

import (
	"fmt"
	"math/rand"
)

var (
	temperature = 15.0
	humidity    = 30.0
)

func GetTemperature() string {
	temperature += rand.Float64()*10 - 4
	return fmt.Sprintf("%f", temperature)
}

func GetHumidity() string {
	humidity += rand.Float64()*10 - 4
	return fmt.Sprintf("%f", humidity)
}
