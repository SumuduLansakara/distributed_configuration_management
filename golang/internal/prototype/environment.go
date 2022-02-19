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

var (
	temperature = 15.0
	humidity    = 30.0
)

func getRandomInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func getTemperature(w http.ResponseWriter, req *http.Request) {
	temperature += getRandomInRange(-1, 1)

	w.Header().Set("Content-Type", "application/json")

	reply := map[string]string{
		"val": fmt.Sprintf("%f", temperature),
	}
	json, err := json.Marshal(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func setTemperature(w http.ResponseWriter, req *http.Request) {
	val, ok := req.URL.Query()["v"]
	if !ok || len(val) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	valF, err := strconv.ParseFloat(val[0], 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	zap.L().Info("temperature value updated", zap.Float64("old", humidity), zap.Float64("new", valF))
	temperature = valF
}

func getHumidity(w http.ResponseWriter, req *http.Request) {
	humidity += getRandomInRange(-1, 1)
	if humidity < 0 {
		humidity = 0
	} else if humidity > 100 {
		humidity = 100
	}
	w.Header().Set("Content-Type", "application/json")

	reply := map[string]string{
		"val": fmt.Sprintf("%f", humidity),
	}
	json, err := json.Marshal(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func setHumidity(w http.ResponseWriter, req *http.Request) {
	val, ok := req.URL.Query()["v"]
	if !ok || len(val) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	valF, err := strconv.ParseFloat(val[0], 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	zap.L().Info("humidity value updated", zap.Float64("old", humidity), zap.Float64("new", valF))
	humidity = valF
}

func StartEnvironment() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc(QueryKeyGetTemperature, getTemperature)
	http.HandleFunc(QueryKeyGetHumidity, getHumidity)
	http.HandleFunc(QueryKeySetTemperature, setTemperature)
	http.HandleFunc(QueryKeySetHumidity, setHumidity)
	go http.ListenAndServe(":3100", nil)
}
