package prototype

import (
	"hash/fnv"
	"math/rand"
	"os"
	"time"

	"go.uber.org/zap"
)

var (
	temperature = 15.0
	humidity    = 30.0
)

func init() {
	// Use a random seed based on the hostname to have a unique seed for each container
	hostname, err := os.Hostname()
	if err != nil {
		zap.L().Panic("failed getting hostname", zap.Error(err))
	}
	h := fnv.New64a()
	h.Write([]byte(hostname))
	hash := h.Sum64()
	rand.Seed(time.Now().UnixNano() * int64(hash))
}

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
