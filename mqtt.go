package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttConfig struct {
	Name              string      `json:"name"`
	UnitOfMeasurement string      `json:"unit_of_measurement"`
	StateTopic        string      `json:"state_topic"`
	ConfigTopic       string      `json:"config_topic"`
	UniqueID          string      `json:"unique_id"`
	Client            mqtt.Client `json:"-"`
}

func (m *mqttConfig) Connect(uri *url.URL) {
	opts := createClientOptions(m.Name, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(10 * time.Second) {
	}

	if err := token.Error(); err != nil {
		log.Fatalln(err)
	}
	m.Client = client
}

func createClientOptions(clientID string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetClientID(clientID)
	return opts
}

// SendConfig send the home assistant auto discovery config to mqtt
func (m *mqttConfig) SendConfig() error {
	mqttPayload, err := preparePayload(m)
	if err != nil {
		return err
	}
	err = sendMessage(m.Client, m.ConfigTopic, mqttPayload)
	return err
}

// Send any data to home assistant
func (m *mqttConfig) Send(data string) error {
	err := sendMessage(m.Client, m.StateTopic, data)
	return err
}

func sendMessage(client mqtt.Client, topic string, data string) error {
	token := client.Publish(topic, 0, false, data)
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

func preparePayload(data interface{}) (string, error) {
	mqttPayload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(mqttPayload), nil
}
