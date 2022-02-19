package prototype

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const (
	EnvPropertyTemperature = "temperature"
	EnvPropertyHumidity    = "humidity"
)

var (
	temperature = 15.0
	humidity    = 30.0
)

func getRandomInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func get(w http.ResponseWriter, req *http.Request) {
	keys, ok := req.URL.Query()["key"]
	if !ok || len(keys) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var value float64
	switch keys[0] {
	case EnvPropertyTemperature:
		temperature += getRandomInRange(-1, 1)
		value = temperature
	case EnvPropertyHumidity:
		humidity += getRandomInRange(-1, 1)
		if humidity < 0 {
			humidity = 0
		} else if humidity > 100 {
			humidity = 100
		}
		value = humidity
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	reply := map[string]string{
		"value": fmt.Sprintf("%f", value),
	}
	json, err := json.Marshal(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func set(w http.ResponseWriter, req *http.Request) {
	values, ok := req.URL.Query()["value"]
	if !ok || len(values) != 1 {
		zap.L().Warn("query parameters missing", zap.String("missing_parameter", "value"), zap.Any("request", req.URL))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	valF, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		zap.L().Warn("unable to parse", zap.Any("val", values[0]))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keys, ok := req.URL.Query()["key"]
	if !ok || len(keys) != 1 {
		zap.L().Warn("query parameters missing", zap.String("missing_parameter", "key"), zap.Any("request", req.URL))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch keys[0] {
	case EnvPropertyTemperature:
		zap.L().Info("temperature value updated", zap.Float64("old", temperature), zap.Float64("new", valF))
		temperature = valF
	case EnvPropertyHumidity:
		zap.L().Info("humidity value updated", zap.Float64("old", humidity), zap.Float64("new", valF))
		humidity = valF
	default:
		zap.L().Warn("invalid key", zap.String("key", keys[0]))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func StartEnvironment() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc(QueryKeyGet, get)
	http.HandleFunc(QueryKeySet, set)
	go http.ListenAndServe(":3100", nil)
}
