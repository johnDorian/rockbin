package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/robfig/cron.v3"
	// "github.com/robfig/cron/v3"
)


var FilePath = "/mnt/data/rockrobo/RoboController.cfg"
var percentage = false;

func main() {

	var mqttServer string
	var sensorName string
	var binFullTime float64
	var unitOfMeasurement string
	flag.StringVar(&mqttServer, "mqtt_server", "mqtt://localhost:1883", "mqtt broker address")
	flag.StringVar(&sensorName, "sensor_name", "vacuumbin", "Name of sensor in Home Assistant")
	flag.Float64Var(&binFullTime, "full_time", 0, "Amount of seconds where the bin will be considered full")
	flag.Parse()

	if binFullTime > 0 {
		percentage = true
	}

	if percentage {
		unitOfMeasurement = "%"
	} else {
		unitOfMeasurement = "min"
	}

	var homeAssConfig = fmt.Sprintf("{\"name\": %[1]q, \"unit_of_measurement\": %[2]q, \"state_topic\": \"homeassistant/sensor/%[1]s/state\", \"unique_id\": %[1]q}", sensorName, unitOfMeasurement)
	var homeAssConfigTopic = fmt.Sprintf("homeassistant/sensor/%s/config", sensorName)
	var stateTopic = fmt.Sprintf("homeassistant/sensor/%s/state", sensorName)

	
	mqttURL, err := url.Parse(mqttServer)
	if err != nil {
		log.Fatalln(err)
	}
	mqttClient := connect("bin", mqttURL)

	// on launch tell home assistant that we exist

	c := cron.New()
	c.AddFunc("@every 0h1m0s", func() {
		token := mqttClient.Publish(homeAssConfigTopic, 0, false, homeAssConfig)
		if token.Error() != nil {
			log.Println(token.Error())
		}
		binTime := getBinValue(binFullTime)

		token = mqttClient.Publish(stateTopic, 0, false, binTime)
		if token.Error() != nil {
			log.Println(token.Error())
		}

	})
	c.Start()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				_ = event
				time.Sleep(time.Second * 1)
				binTime := getBinValue(binFullTime)

				token := mqttClient.Publish(stateTopic, 0, false, binTime)
				if token.Error() != nil {
					log.Println(token.Error())
				}
			case err := <-watcher.Errors:
				log.Fatalln(err)

			}
		}
	}()

	if err := watcher.Add(FilePath); err != nil {
		log.Fatalln(err)
	}

	<-done

}



func getBinValue(binFullTime float64) string {

	//FilePath := "RobotController.cfg"
	file, err := os.Open(FilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := ""
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Contains(line, "bin_in_time") {
			line = strings.Split(line, "=")[1]
			line = strings.Trim(line, " ;")
			break
		}
	}
	file.Close()
	BinTime, err := strconv.ParseFloat(line, 32)
	if percentage {
		binCapacity := BinTime / binFullTime * 100.
		return fmt.Sprintf("%.2f", binCapacity)
	} else {
		binCapacity := BinTime
		return fmt.Sprintf("%d", int(binCapacity/60))
	}
}

func connect(clientID string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientID, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(10 * time.Second) {
	}

	if err := token.Error(); err != nil {
		log.Fatalln(err)
	}

	return client
}

func createClientOptions(clientID string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetClientID(clientID)
	return opts
}
