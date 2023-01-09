package cmd

import (
	"errors"
	"os"

	"github.com/graytonio/configarr/internal/api"
	"github.com/graytonio/configarr/internal/config"
	"github.com/graytonio/configarr/internal/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	forceFlag     bool
	VersionString string = "dev"
	rootCmd              = &cobra.Command{
		SilenceUsage: true,
		Use:          "configarr",
		Version:      VersionString,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, errs := config.ParseConfig()

			if len(errs) > 0 && !forceFlag {
				return errors.New("failed to initialize services")
			}

			for _, service := range config.Services {
				err := api.ApplyServiceConfig(&service)
				if err != nil {
					log.Logger.WithFields(logrus.Fields{
						"error":     err.Error(),
						"component": service.Name,
					}).Error("Could Not Configure Service")
					log.Logger.WithField("forceFlag", forceFlag).Debug("Force Flag")
					if !forceFlag {
						return err
					}
				}
			}

			return nil
		},
	}
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file default: ./config.yaml")
	rootCmd.PersistentFlags().BoolVarP(&forceFlag, "force", "f", false, "Continue trying to configure services if one fails")
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
