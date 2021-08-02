package main

import (
	"flag"
	"fmt"
	_ "github.com/iAnatoly/telegraf-input-dht_sensor/plugins/inputs/dht_sensor"
	"github.com/influxdata/telegraf/plugins/common/shim"
	"os"
	"time"
)

var pollInterval = flag.Duration("poll_interval", 1*time.Second, "how often to send metrics")
var pollIntervalDisabled = flag.Bool("poll_interval_disabled", false, "set to true to disable polling. You want to use this when you are sending metrics on your own schedule")
var configFile = flag.String("config", "", "path to the config file for this plugin")
var err error

func main() {
	// parse command line options
	flag.Parse()
	if *pollIntervalDisabled {
		*pollInterval = shim.PollIntervalDisabled
	}

	// create the shim. This is what will run the plugin.
	shim := shim.New()

	// If no config is specified, all imported plugins are loaded.
	err = shim.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err loading input: %s\n", err)
		os.Exit(1)
	}

	// run a single plugin until stdin closes or we receive a termination signal
	if err := shim.Run(*pollInterval); err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		os.Exit(1)
	}
}
