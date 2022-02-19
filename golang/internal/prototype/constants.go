package prototype

const (
	KindSensor     = "sensor"
	KindDisplay    = "display"
	KindController = "controller"
	KindActuator   = "actuator"
)

const (
	ParamSensorType  = "sensorType"
	ParamTemperature = "temperature"
	ParamHumidity    = "humidity"

	ParamActuatorType = "actuatorType"

	ParamACState = "acState"
)

const (
	ValueSensorTypeTemperatureSensor = "temperature"
	ValueSensorTypeHumiditySensor    = "humidity"

	ValueActuatorTypeAirConditioner = "airconditioner"

	ValueACStateActive   = "active"
	ValueACStateInactive = "inactive"
)

const (
	QueryKeyGetTemperature = "/gettemperature"
	QueryKeySetTemperature = "/settemperature"
	QueryKeyGetHumidity    = "/gethumidity"
	QueryKeySetHumidity    = "/sethumidity"
)
