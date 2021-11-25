package mqtt

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttConfig struct {
	Name              string      `json:"name"`
	UnitOfMeasurement string      `json:"unit_of_measurement"`
	StateTopic        string      `json:"state_topic"`
	UniqueID          string      `json:"unique_id"`
	MaxConnectionTime int         `json:"-"`
	Client            mqtt.Client `json:"-"`
	ConfigTopic       string      `json:"-"`
	Server            string      `json:"-"`
	Username          string      `json:"-"`
	Password          string      `json:"-"`
}

func (m *MqttConfig) ConnectWithBackoff() error {
	connectionBackoff := backoff.NewExponentialBackOff()
	connectionBackoff.InitialInterval = 1 * time.Second
	connectionBackoff.MaxElapsedTime = time.Duration(m.MaxConnectionTime) * time.Second
	err := backoff.RetryNotifyWithTimer(m.Connect,
		connectionBackoff,
		func(e error, d time.Duration) {
			log.Debug("mqtt connection attempt failed. Trying again in ", d.String())
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("mqtt connection to: %s failed after trying for %s seconds", m.Server, connectionBackoff.MaxElapsedTime)
	}
	return nil
}

func (m *MqttConfig) Connect() error {
	mqttURL, err := url.Parse(m.Server)
	if err != nil {
		return err
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", mqttURL.Host))
	opts.SetClientID(m.UniqueID)
	if len(m.Username) > 0 {
		opts.SetUsername(m.Username)
	}
	if len(m.Password) > 0 {
		opts.SetPassword(m.Password)
	}

	client := mqtt.NewClient(opts)

	token := client.Connect()
	for !token.WaitTimeout(2 * time.Second) {
	}

	if err := token.Error(); err != nil {
		return err
	}
	if client.IsConnected() {
		log.WithFields(log.Fields{"mqtt_broker": mqttURL.Host}).Info("Connected to mqtt broker")
	}

	m.Client = client
	return nil
}

// SendConfig send the home assistant auto discovery config to mqtt
func (m *MqttConfig) SendConfig() error {
	mqttPayload, err := preparePayload(m)
	if err != nil {
		return err
	}
	log.Debug("Sending mqtt config")
	err = sendMessage(m.Client, m.ConfigTopic, mqttPayload, true)
	return err
}

// Send any data to home assistant
func (m *MqttConfig) Send(data string) error {
	log.Debug("Sending mqtt message")
	err := sendMessage(m.Client, m.StateTopic, data, false)
	return err
}

func sendMessage(client mqtt.Client, topic string, data string, retain bool) error {
	token := client.Publish(topic, 0, retain, data)
	if token.Error() != nil {
		log.Fatalln(token.Error())
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
