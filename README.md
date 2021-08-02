# DHT Sensor plugin

## Summary

This plugin collects temperature and humidity from DHTXX sensors (usually connected to Raspberry Pi or similar host). 

## Configuration

Sample config (see plugin.conf in the repo):
```toml
[[inputs.dht_sensor]]

    # type of the sensor you are using (check your invoice if not sure).
    # possible values: "DHT11", "DHT12", "DHT22"
    sensor = "DHT22"

    # which data pin are you connectinmg to? 
    # Check your wiring if you are not sure
    data_pin = 4

    # see https://github.com/d2r2/go-dht for the meaning of this one.
    # TL/DR: play with this parameter    
    boost_gpio_performance = true

    # number of retries in case of a failure
    number_of_retries = 3
    
    # also retry if a sensor returns data outside the sanity range.
    # temperature is specified in degrees celsius.
    retry_humidity_above = 100
    retry_temperature_above = 80 

```

## Installation

* Clone the repo
```bash
git clone 
```
* Build the "dht_sensor" binary

```bash
$ go build -o dht_sensor cmd/main.go
# or, if you need to cross-compile for arm, because you are collecting the data from Raspberry Pi probes:
$ env GOOS=linux GOARCH=arm GOARM=7 go build -o dht_sensor.arm7 cmd/main.go
```
* Edit the config
* Copy the binary and the config to an appropriate location
```bash
$ sudo cp plugin.config /etc/telegraf/telegraf-dht_sensor.config
$ sudo cp dht_sensor /usr/lib/telegraf/plugins/
```
* You should be able to call this from telegraf now using execd
```
[[inputs.execd]]
  command = ["/usr/lib/telegraf/plugins/dht_sensor", "-config", "/etc/telegraf/telegraf-dht_sensor.config" ]
  signal = "none"
```
## Credits
* This self-contained plugin is based on the documentations of [Execd Go Shim](https://github.com/influxdata/telegraf/blob/master/plugins/common/shim).
* This plugin uses [Go-DHT](https://github.com/d2r2/go-dht) library, Most of the heavy lifting is done there. Huge kudos to the authors.
