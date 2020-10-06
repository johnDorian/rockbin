package main

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
