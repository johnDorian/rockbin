/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/johnDorian/rockbin/status"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of the running service",
	Long:  `This sub command will query the http endpoint as defined in the config for the current status`,
	Run: func(cmd *cobra.Command, args []string) {
		endPoint := fmt.Sprintf("http://%s:%s/status", viper.GetString("status_address"), viper.GetString("status_port"))
		resp, err := http.Get(endPoint)
		if err != nil {
			log.Fatalln("Web server not found", err)
		}
		defer resp.Body.Close()

		stats := status.Data{}
		json.NewDecoder(resp.Body).Decode(&stats)

		prettyOutput, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			log.Fatalln("Error formatting response", err)
		}
		fmt.Println(string(prettyOutput))

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
