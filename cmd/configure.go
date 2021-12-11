package cmd

import (
	"log"

	"github.com/johnDorian/rockbin/configure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var servicePath string

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the config and startup script for the server",
	Long: `This is an optional sub-command which provides the ability
to interactively create a config file and add the startup script to 
the correct location.`,
	Run: func(cmd *cobra.Command, args []string) {
		prompter, err := configure.NewPrompter(serveCmd.Flags(), viper.GetViper().ConfigFileUsed(), servicePath)
		if err != nil {
			log.Fatalln(err)
		}
		err = prompter.PromptUser()
		if err != nil {
			log.Fatalln(err)
		}

		err = prompter.WriteOutTemplate("config", prompter.Responses)
		if err != nil {
			log.Fatalln(err)
		}
		err = prompter.WriteOutTemplate("service", viper.GetViper().ConfigFileUsed())
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().StringVar(&servicePath, "service_path", "", "path to store the service file")

}
