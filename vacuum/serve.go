package vacuum

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/johnDorian/rockbin/mqtt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func Serve(bin Bin, mqttClient mqtt.MqttConfig) {
	// on launch tell home assistant that we exist
	mqttClient.SendConfig()

	// every minute send everything to the mqtt broker
	c := cron.New()
	c.AddFunc("@every 0h1m0s", func() {
		log.Debug("Running cron job")
		mqttClient.SendConfig()
		bin.Update()
		mqttClient.Send(bin.Value)
	})
	c.Start()

	// Setup a file watcher to get instance updates on file changes
	log.Debug("Setting up file watcher for: ", bin.FilePath)
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

				bin.Update()
				log.Debug("Current bin value is:", bin.Value)
				mqttClient.Send(bin.Value)
			case err := <-watcher.Errors:
				log.Fatalln(err)

			}
		}
	}()

	if err := watcher.Add(bin.FilePath); err != nil {
		log.Fatalln(err)
	}

	<-done

}
