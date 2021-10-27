package configure

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/manifoldco/promptui"
	"github.com/spf13/pflag"
)

type CLIVariables struct {
	Variables   []cmdOption
	Responses   map[string]string
	ServiceFile string
	ConfigFile  string
}

// Taken from https://github.com/spf13/cobra/blob/master/doc/yaml_docs.go
type cmdOption struct {
	Name         string
	Shorthand    string
	DefaultValue string
	Usage        string
}

//go:embed templates/*
var templateFiles embed.FS

// Modified from https://github.com/spf13/cobra/blob/master/doc/yaml_docs.go
func genFlagResult(flags *pflag.FlagSet) []cmdOption {
	var result []cmdOption

	flags.VisitAll(func(flag *pflag.Flag) {
		// Todo, when we mark a shorthand is deprecated, but specify an empty message.
		// The flag.ShorthandDeprecated is empty as the shorthand is deprecated.
		// Using len(flag.ShorthandDeprecated) > 0 can't handle this, others are ok.
		if !(len(flag.ShorthandDeprecated) > 0) && len(flag.Shorthand) > 0 {
			opt := cmdOption{
				flag.Name,
				flag.Shorthand,
				flag.DefValue,
				flag.Usage,
			}
			result = append(result, opt)
		} else {
			opt := cmdOption{
				Name:         flag.Name,
				DefaultValue: flag.DefValue,
				Usage:        flag.Usage,
			}
			result = append(result, opt)
		}
	})

	return result
}

func NewPrompter(flags *pflag.FlagSet, configFile, serviceFilePath string) (*CLIVariables, error) {
	setup := CLIVariables{
		Variables: genFlagResult(flags),
	}

	if serviceFilePath == "" {
		serviceFilePath = "/etc/init/"
	}

	serviceFile := "rockbin.conf"

	//determine if upstart or not
	_, err := os.Stat("/sbin/initctl")
	if err != nil {
		serviceFile = "S12Rockbin"
	}

	setup.ConfigFile = configFile
	setup.ServiceFile = path.Join(serviceFilePath, serviceFile)

	return &setup, nil
}

func (p *CLIVariables) PromptUser() error {

	p.Responses = make(map[string]string)
	for _, arg := range p.Variables {
		prompt := promptui.Prompt{
			Label:   arg.Usage,
			Default: arg.DefaultValue,
		}

		answer, err := prompt.Run()

		if err != nil {
			return err
		}

		p.Responses[arg.Name] = answer
	}
	return nil
}

func (p *CLIVariables) WriteOutTemplate(file string, data interface{}) error {

	var (
		outputFile string
		tmplFile   string
	)

	switch file {
	case "config":
		tmplFile = "templates/config.yaml.tmpl"
		outputFile = p.ConfigFile

	case "service":
		tmplFile = path.Join("templates", fmt.Sprintf("%s.%s", path.Base(p.ServiceFile), "tmpl"))
		outputFile = p.ServiceFile
	}
	t, err := template.ParseFS(templateFiles, tmplFile)
	if err != nil {
		return err
	}

	os.MkdirAll(path.Dir(outputFile), 0755)

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, data)

	fmt.Printf("Writing %s file to: %s\n", file, outputFile)

	return err

}
