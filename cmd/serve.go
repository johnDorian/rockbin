package cmd

import (
	"fmt"
	"net/url"

	"github.com/johnDorian/rockbin/mqtt"
	"github.com/johnDorian/rockbin/vacuum"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the daemon",
	Long: `Start the daemon for sending the bin value to 
the mqtt server.`,

	Run: func(cmd *cobra.Command, args []string) {

		setUpLogger(viper.GetString("log_level"))

		bin := vacuum.Bin{
			FilePath: viper.GetString("file_path"),
			Capacity: viper.GetFloat64("full_time"),
			Unit:     viper.GetString("measurement_unit"),
		}

		mqttURL, err := url.Parse(viper.GetString("mqtt_server"))
		if err != nil {
			log.Fatalln(err)
		}

		mqttClient := mqtt.MqttConfig{
			Name:              viper.GetString("sensor_name"),
			UnitOfMeasurement: viper.GetString("measurement_unit"),
			StateTopic:        fmt.Sprintf(viper.GetString("mqtt_state_topic"), viper.GetString("sensor_name")),
			ConfigTopic:       fmt.Sprintf("homeassistant/sensor/%v/config", viper.GetString("sensor_name")),
			UniqueID:          viper.GetString("sensor_name"),
		}
		mqttClient.Connect(mqttURL, viper.GetString("mqtt_user"), viper.GetString("mqtt_password"))
		vacuum.Serve(bin, mqttClient)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("mqtt_server", "mqtt://localhost:1883", "Address of the mqtt server")
	serveCmd.Flags().String("mqtt_user", "", "Username for the mqtt server")
	serveCmd.Flags().String("mqtt_password", "", "Password for the mqtt server")
	serveCmd.Flags().String("mqtt_state_topic", "homeassistant/sensor/%v/state", "State topic (%v is replaced with the sensor_name value)")
	serveCmd.Flags().String("sensor_name", "vacuumbin", "Name of sensor in Home Assistant")
	serveCmd.Flags().Float64("full_time", 2400., "Amount of seconds where the bin will be considered full")
	serveCmd.Flags().String("measurement_unit", "%", "In what unit should the measurement be sent (%, sec, min)")
	serveCmd.Flags().String("file_path", "/mnt/data/rockrobo/RoboController.cfg", "File path of RoboController.cfg")
	serveCmd.Flags().String("log_level", "Fatal", "Level of logging (trace, debug, info, warn, error, fatal, panic).")

}

func setUpLogger(level string) {
	loglevel, _ := log.ParseLevel(level)
	log.SetLevel(loglevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.Info("Starting rockbin service")
	log.WithFields(log.Fields{"loglevel": log.GetLevel()}).Debug("Setup logger with log level")
}
