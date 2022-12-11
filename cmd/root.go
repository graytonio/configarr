package cmd

import (
	log "github.com/sirupsen/logrus"
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
			// service, err := api.NewService("localhost", 7878)
			// if err != nil {
			// 	return err
			// }

			// log.Infof("API Key: %s\n", service.ApiKey)

			log.Info(viper.AllKeys())

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
		log.Fatalf("Error parsing verbose flag: %s\n", err.Error())
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s\n", viper.ConfigFileUsed())
	} else {
		log.Errorf("Error reading config file: %s", err.Error())
	}
}
