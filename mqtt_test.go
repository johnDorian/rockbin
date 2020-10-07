package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

func TestcreateClientOptions(t *testing.T) {
	assert := assert.New(t)

	var testData = []struct {
		clientID string
		hostname string
		username string
		password string
	}{
		{"something", "http://test.com", "", ""},
		{"hello-192!", "http://test.com", "", ""},
		{"something", "http://test.com", "hello", ""},
		{"something", "http://test.com", "hello", "world"},
		{"something", "http://test.com", "h", "w"},
	}

	for _, test := range testData {
		expected := mqtt.NewClientOptions()
		uri, _ := url.Parse(test.hostname)
		expected.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
		expected.SetClientID(test.clientID)
		if len(test.username) > 0 {
			expected.SetUsername(test.username)
		}
		if len(test.password) > 0 {
			expected.SetPassword(test.password)
		}

		os.Setenv("MQTT_USERNAME", test.username)
		os.Setenv("MQTT_PASSWORD", test.password)
		opts := createClientOptions(test.clientID, uri)
		assert.Equal(expected, opts)
	}

}

func TestConnect(t *testing.T) {
	assert := assert.New(t)
	var testData = []struct {
		username string
		password string
	}{
		{"user", "pass"},
		{"user1", "!$%&/()?#+*12345"},
		{"user2", `hello"world`},
	}
	resource, pool := spinUpMQTT()
	for _, up := range testData {
		config := mqttConfig{Name: "hello", UnitOfMeasurement: "hello", StateTopic: "hello", ConfigTopic: "hello", UniqueID: "hello"}
		uri, _ := url.Parse(fmt.Sprintf("mqtt://127.0.0.1:%v", resource.GetPort("1883/tcp")))
		os.Setenv("MQTT_USERNAME", up.username)
		os.Setenv("MQTT_PASSWORD", up.password)
		config.Connect(uri)
		assert.True(config.Client.IsConnected())

	}

	pool.Purge(resource)

}

func spinUpMQTT() (*dockertest.Resource, *dockertest.Pool) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	dir, _ := os.Getwd()

	options := &dockertest.RunOptions{
		Repository: "eclipse-mosquitto",
		Tag:        "latest",

		ExposedPorts: []string{"1883"},
		Mounts: []string{fmt.Sprintf("%v:/mosquitto/config/mosquitto.conf", path.Join(dir, "tests/mosquitto.conf")),
			fmt.Sprintf("%v:/password.txt", path.Join(dir, "tests/password.txt"))},
	}

	resource, err := pool.RunWithOptions(options)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	return resource, pool
}
