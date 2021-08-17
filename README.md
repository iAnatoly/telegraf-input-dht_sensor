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
```

NOTE: if you are building for Raspbery PI, you need to either:
1. Build it on the Raspberry Pi device itself, or
2. Cross-compile. 

Usually (2) should be quite simple, but go-dht library has a portion of the code implemented in C, so you will need to install a cross compiler for ARM, and cross-compile using the following:
``` bash
# compile the module for amd64.
$ go build -o net_irtt.amd64 cmd/main.go

# install the debian package for cross-compilation
$ sudo apt install gcc-10-arm-linux-gnueabi

# compile for ARM64 (RPi4)
$ env GOOS=linux GOARCH=arm64 go build -o net_irtt.arm64 cmd/main.go

# compile for ARMv7l (RPI3b)
$ env CC=arm-linux-gnueabi-gcc-10 CGO_ENABLED=1  GOOS=linux GOARCH=arm GOARM=7 go build -o dht_sensor.armv7l cmd/main.go

# compile for ARMv6l (RPI Zero)
$ env CC=arm-linux-gnueabi-gcc-10 CGO_ENABLED=1  GOOS=linux GOARCH=arm GOARM=6 go build -o dht_sensor.armv6l cmd/main.go

```
(notice the extra CC and CGO_ENABLED variables required for ARM6/7 in addition to regular Go cross-compile flags).
It is probably easier to compile directly on RPi, but some peopel do not like to install DEV dependencies in production. 

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
* If you get an error `Error during call C.dial_DHTxx_and_read(): failed to open GPIO export for writing`, this means that the user you are running telegraf under does not have permissions to open `/sys/class/gpio/export`. To fix this, simply add your telegraf user into gpio group, and restart telegraf:
```bash
$ systemctl status telegraf 
[...]
Aug 02 14:43:59 pi telegraf[2419]: 2021-08-02T21:43:59Z E! [inputs.execd] stderr: "failed to gather metrics: Error during call C.dial_DHTxx_and_read(): failed to open GPIO export for writing"
$ ls -l /sys/class/gpio/export
-rwxrwx--- 1 root gpio 4096 Jul 31 20:17 /sys/class/gpio/export
$ sudo usermod -a -G gpio telegraf
$ sudo systemctl restart telegraf

```

## Credits
* This self-contained plugin is based on the documentations of [Execd Go Shim](https://github.com/influxdata/telegraf/blob/master/plugins/common/shim).
* This plugin uses a fork of [Go-DHT](https://github.com/d2r2/go-dht) library. Most of the heavy lifting is done there. Huge kudos to the authors.
