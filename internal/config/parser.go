package config

import (
	"github.com/graytonio/configarr/internal/log"
	"github.com/spf13/viper"
)

func ParseConfig() *Config {
	log.Logger.WithField("component", "config").Info("Reading Config File")
	var configData = Config{}
	err := viper.Unmarshal(&configData)
	if err != nil {
		log.Logger.WithField("error", err.Error()).Fatal("Failed to Unmarshal Config")
	}

	// Attempt to pull init data from each service do not continue if it fails
	var goodServices []Service
	for idx := range configData.Services {
		err := configData.Services[idx].InitService()
		if err != nil {
			log.Logger.Errorf("Failed to setup service %s: %s", configData.Services[idx].Name, err.Error())
			continue
		}
		goodServices = append(goodServices, configData.Services[idx])
		log.Logger.WithField("component", "config").Debug(configData.Services[idx])
	}

	configData.Services = goodServices
	log.Logger.Debug(configData)

	populateProwlarrApplicationData(&configData)

	return &configData
}

func populateProwlarrApplicationData(config *Config) {
	for srv_idx, srv := range config.Services {
		if len(srv.Config.Applications) == 0 {
			continue
		}

		for app_idx, app := range srv.Config.Applications {

			config.Services[srv_idx].Config.Applications[app_idx].Fields = append(config.Services[srv_idx].Config.Applications[app_idx].Fields, ResourceField{Name: "prowlarrUrl", Value: srv.ServiceAddress})

			// Find refernced app and set fields
			for _, service := range config.Services {
				if app.Name != service.Name {
					continue
				}
				config.Services[srv_idx].Config.Applications[app_idx].Fields = append(config.Services[srv_idx].Config.Applications[app_idx].Fields, ResourceField{Name: "apiKey", Value: service.ApiKey})
				config.Services[srv_idx].Config.Applications[app_idx].Fields = append(config.Services[srv_idx].Config.Applications[app_idx].Fields, ResourceField{Name: "baseUrl", Value: service.ServiceAddress})
			}
		}
	}
}
