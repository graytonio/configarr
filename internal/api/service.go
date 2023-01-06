package api

import (
	"encoding/json"
	"fmt"

	"github.com/graytonio/configarr/internal/config"
	"github.com/graytonio/configarr/internal/log"
	"github.com/graytonio/configarr/internal/rest"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func ApplyServiceConfig(service *config.Service) error {
	if service.IsEmpty() {
		return nil
	}

	logger := log.Logger.WithField("component", service.Name)
	logger.Info("Applying Configuration")
	logger.WithField("service", service).Debug("Service Object")

	err := applyConfigResource("rootfolder", service, service.Config.RootFolders, logger)
	if err != nil {
		return err
	}

	err = applyConfigResource("downloadclient", service, service.Config.DownloadClients, logger)
	if err != nil {
		return err
	}

	err = applyConfigResource("applications", service, service.Config.Applications, logger)
	if err != nil {
		return err
	}

	return nil
}

func applyConfigResource[R config.ResourceConfig](resourceName string, service *config.Service, configData []R, logger *logrus.Entry) error {
	if len(configData) == 0 {
		logger.Infof("No %s to Configure", resourceName)
		return nil
	}

	logger.Infof("Applying %s Configuration", resourceName)

	// Get current configuration of resource
	currentResourceApiResponse, err := rest.GetResource(service.ApiAddress, service.ApiKey, resourceName)
	if err != nil {
		return err
	}

	var currentResources []R
	err = json.Unmarshal(currentResourceApiResponse, &currentResources)
	if err != nil {
		return nil
	}

	logger.WithField(resourceName, currentResources).Debug("Current Configuration")

	// Delete any resources not found in configuration file
	for _, resource := range currentResources {
		if !slices.ContainsFunc(configData, func(f R) bool { return f.GetName() == resource.GetName() }) {
			logger.Infof("Deleting %s: %v", resourceName, resource.GetName())
			err = rest.DeleteResource(service.ApiAddress, service.ApiKey, resourceName, fmt.Sprint(resource.GetID()))
			if err != nil {
				return err
			}
		}
	}

	// Create any resources found in the configuration file that don't already exist
	for _, resource := range configData {
		if !slices.ContainsFunc(currentResources, func(f R) bool { return f.GetName() == resource.GetName() }) {
			logger.Infof("Creating %s: %v", resourceName, resource.GetName())
			err = rest.PostResource(service.ApiAddress, service.ApiKey, resourceName, resource.GetCreateData())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
