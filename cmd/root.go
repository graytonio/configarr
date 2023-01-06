package cmd

import (
	"github.com/graytonio/configarr/internal/api"
	"github.com/graytonio/configarr/internal/config"
	"github.com/graytonio/configarr/internal/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	VersionString string = "dev"
	rootCmd              = &cobra.Command{
		Use:     "configarr",
		Version: VersionString,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := config.ParseConfig()

			for _, service := range config.Services {
				err := api.ApplyServiceConfig(&service)
				if err != nil {
					log.Logger.WithFields(logrus.Fields{
						"error":     err.Error(),
						"component": service.Name,
					}).Error("Could Not Configure Service")
				}
			}

			return nil
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file default: ./config.yaml")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose logging")
}

func initConfig() {
	verbose, err := rootCmd.Flags().GetBool("verbose")
	if err != nil {
		log.Logger.WithField("error", err.Error()).Fatal("Error parsing verbose flag")
	}

	if verbose {
		log.Logger.SetLevel(logrus.DebugLevel)
	}

	log.Logger.SetFormatter(&log.Formatter{
		HideKeys:    !verbose,
		FieldsOrder: []string{"component", "error"},
	})

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Logger.WithField("config_file", viper.ConfigFileUsed()).Debug("Parsed Config File")
	} else {
		log.Logger.WithField("error", err.Error()).Error("Error reading config file")
	}
}
