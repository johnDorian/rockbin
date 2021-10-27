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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
