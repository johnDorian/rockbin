package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

func config() (Bin, mqttConfig) {
	var mqttServer string
	var mqttUser string
	var mqttPassword string
	var mqttStateTopic string
	var sensorName string
	var binFullTime float64
	var unitOfMeasurement string
	var FilePath string
	var LoggingLevel string
	flag.StringVar(&mqttServer, "mqtt_server", "mqtt://localhost:1883", "mqtt broker address")
	flag.StringVar(&mqttUser, "mqtt_user", "", "mqtt user")
	flag.StringVar(&mqttPassword, "mqtt_password", "", "mqtt password")
	flag.StringVar(&mqttStateTopic, "mqtt_state_topic", "homeassistant/sensor/%v/state", "State topic (%v is replaced with the sensor_name value)")
	flag.StringVar(&sensorName, "sensor_name", "vacuumbin", "Name of sensor in Home Assistant")
	flag.Float64Var(&binFullTime, "full_time", 2400., "Amount of seconds where the bin will be considered full")
	flag.StringVar(&unitOfMeasurement, "measurement_unit", "%", "In what unit should the measurement be sent (%, sec, min)")
	flag.StringVar(&FilePath, "file_path", "/mnt/data/rockrobo/RoboController.cfg", "file path of RoboController.cfg")
	flag.StringVar(&LoggingLevel, "log_level", "Fatal", "Level of logging (trace, debug, info, warn, error, fatal, panic).")
	flag.Parse()

	setUpLogger(LoggingLevel)
	printVersion()

	bin := Bin{
		FilePath: FilePath,
		Capacity: binFullTime,
		Unit:     unitOfMeasurement,
	}

	mqttURL, err := url.Parse(mqttServer)
	if err != nil {
		log.Fatalln(err)
	}

	mqttClient := mqttConfig{
		Name:              sensorName,
		UnitOfMeasurement: unitOfMeasurement,
		StateTopic:        fmt.Sprintf(mqttStateTopic, sensorName),
		ConfigTopic:       fmt.Sprintf("homeassistant/sensor/%v/config", sensorName),
		UniqueID:          sensorName,
	}
	mqttClient.Connect(mqttURL, mqttUser, mqttPassword)
	return bin, mqttClient
}

func printVersion() {
	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			fmt.Println(Version)
			os.Exit(0)
		}
	}
}

func setUpLogger(level string) {
	loglevel, _ := log.ParseLevel(level)
	log.SetLevel(loglevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.Info("Starting rockbin service")
	log.WithFields(log.Fields{"loglevel": log.GetLevel()}).Debug("Setup logger with log level")
}
