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

		// os.Setenv("MQTT_USERNAME", test.username)
		// os.Setenv("MQTT_PASSWORD", test.password)
		opts := createClientOptions(test.clientID, uri, test.username, test.password)
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
	log.Println(resource.GetPort("1883/tcp"))
	for _, up := range testData {
		config := mqttConfig{Name: "hello", UnitOfMeasurement: "hello", StateTopic: "hello", ConfigTopic: "hello", UniqueID: "hello"}
		uri, _ := url.Parse(fmt.Sprintf("mqtt://localhost:%v", resource.GetPort("1883/tcp")))
		//uri, _ := url.Parse("mqtt://localhost:1883")
		// os.Setenv("MQTT_USERNAME", up.username)
		// os.Setenv("MQTT_PASSWORD", up.password)
		log.Println("connecting")
		config.Connect(uri, up.username, up.password)
		assert.True(config.Client.IsConnected())

	}

	pool.Purge(resource)

}

func TestpreparePayload(t *testing.T) {
	assert := assert.New(t)
	var TestData = []struct {
		item1       string
		item2       string
		expected    string
		expectError bool
	}{
		{"item1", "item2", `{"item1":"item1","item2":"item2"}`, false},
	}

	for _, test := range TestData {
		data := struct {
			Item1 string `json:"item1"`
			Item2 string `json:"item2"`
		}{
			test.item1, test.item2,
		}
		payload, err := preparePayload(data)
		assert.Equal(payload, test.expected)
		if test.expectError {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}

}

func spinUpMQTT() (*dockertest.Resource, *dockertest.Pool) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	dir, _ := os.Getwd()
	log.Println("Configuring options")
	options := &dockertest.RunOptions{
		Repository:   "eclipse-mosquitto",
		Tag:          "latest",
		Name:         "mosquitto",
		ExposedPorts: []string{"1883", "9001"},
		Mounts: []string{fmt.Sprintf("%v:/mosquitto/config/mosquitto.conf", path.Join(dir, "tests/mosquitto.conf")),
			fmt.Sprintf("%v:/password.txt", path.Join(dir, "tests/password.txt"))},
	}

	resource, err := pool.RunWithOptions(options)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	log.Println("Up and running")
	return resource, pool
}
