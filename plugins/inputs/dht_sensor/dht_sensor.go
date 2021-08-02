package dht_sensor

import (
	"fmt"
	"strings"
	// sensor imports:
	logger "github.com/d2r2/go-logger"
	"github.com/iAnatoly/go-dht" // use my own fork of 	"github.com/d2r2/go-dht"
	// telegraf imports:
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

const measurement = "dht_sensors"

type DhtSensor struct {
	Sensor                string
	DataPin               int     `toml:"data_pin"`
	Retries               int     `toml:"number_of_retries"`
	RetryHumidityAbove    float32 `toml:"retry_humidity_above"`
	RetryTemperatureAbove float32 `toml:"retry_temperature_above"`
}

func init() {
	logger.ChangePackageLogLevel("dht", logger.PanicLevel)
	inputs.Add("dht_sensor", func() telegraf.Input {
		return &DhtSensor{
			DataPin: 4,
		}
	})
}

func (s *DhtSensor) Description() string {
	return "collects temperature and humidity from DHTXX sensors"
}

// SampleConfig returns sample configuration options.
func (s *DhtSensor) SampleConfig() string {
	return `
    # type of the sensor you are using (check your invoice if not sure).
    # possible values: "DHT11", "DHT12", "DHT22"
    sensor = "DHT22"

    # which data pin are you connectinmg to? 
    # Check your wiring if you are not sure
    data_pin = 4

    # number of retries in case of a failure
    number_of_retries = 3
    
    # also retry if a sensor returns data outside the sanity range.
    # temperature is specified in degrees celsius.
    retry_humidity_above = 100.0
    retry_temperature_above = 80.0 
`
}

func atoSensorType(s string) (dht.SensorType, error) {
	switch strings.ToUpper(s) {
	case "DHT11":
		return dht.DHT11, nil
	case "DHT12":
		return dht.DHT12, nil
	case "DHT22":
		return dht.DHT22, nil
	}

	return -1, fmt.Errorf("Unknown sensor type: %s", s)
}

// Gather is the interface for passing metrics to telegraf
func (n *DhtSensor) Gather(acc telegraf.Accumulator) error {

	sensorType, err := atoSensorType(n.Sensor)
	if err != nil {
		return err
	}

	for i := 0; i < n.Retries; i++ {
		temperature, humidity, _, err := dht.ReadDHTxxWithRetry(sensorType, n.DataPin, false, n.Retries)
		if err != nil {
			return err
		}

		// retry if temperature or humidity are unreasonable
		if temperature > n.RetryTemperatureAbove || humidity > n.RetryHumidityAbove {
			continue
		}

		fields := map[string]interface{}{
			"temp":     temperature,
			"humidity": humidity,
		}

		tags := map[string]string{}

		acc.AddFields(measurement, fields, tags)
		return nil
	}
	return fmt.Errorf("Could not get a reasonable reading after %d retries", n.Retries)
}
