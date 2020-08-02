package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

func config() (Bin, mqttConfig) {
	var mqttServer string
	var sensorName string
	var binFullTime float64
	var unitOfMeasurement string
	var FilePath string
	flag.StringVar(&mqttServer, "mqtt_server", "mqtt://localhost:1883", "mqtt broker address")
	flag.StringVar(&sensorName, "sensor_name", "vacuumbin", "Name of sensor in Home Assistant")
	flag.Float64Var(&binFullTime, "full_time", 2400., "Amount of seconds where the bin will be considered full")
	flag.StringVar(&unitOfMeasurement, "measurement_unit", "%", "In what unit should the measurement be sent (%, sec, min)")
	flag.StringVar(&FilePath, "file_path", "/mnt/data/rockrobo/RoboController.cfg", "file path of RoboController.cfg")
	flag.Parse()

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
		StateTopic:        fmt.Sprintf("homeassistant/sensor/%v/state", sensorName),
		ConfigTopic:       fmt.Sprintf("homeassistant/sensor/%v/config", sensorName),
		UniqueID:          sensorName,
	}
	mqttClient.Connect(mqttURL)
	return bin, mqttClient
}

func printVersion() {
	if os.Args[1] == "version" {
		fmt.Println(Version)
		os.Exit(0)
	}
}
