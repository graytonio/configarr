package config

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

func NewService(name string, address string, port uint16, config ServiceConfig) (*Service, error) {
	var service = Service{
		Name:    name,
		Address: address,
		Port:    port,
		Config:  config,
	}

	err := service.InitService()
	if err != nil {
		return nil, err
	}

	return &service, err
}

func (s *Service) InitService() error {
	// Format url correctly
	serviceAddress := s.Address
	if !strings.HasPrefix(serviceAddress, "http") {
		serviceAddress = fmt.Sprintf("http://%s", serviceAddress)
	}
	serviceAddress = fmt.Sprintf("%s:%d", serviceAddress, s.Port)
	s.ServiceAddress = serviceAddress

	initializeAddress, err := url.JoinPath(serviceAddress, "initialize.js")
	if err != nil {
		return err
	}

	// Fetch initialization variables
	resp, err := http.Get(initializeAddress)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Extract json data from result
	jsonData := strings.Trim(strings.Split(string(body), "=")[1], " ")

	err = yaml.Unmarshal([]byte(jsonData), &s)
	if err != nil {
		return err
	}

	apiAddress, err := url.JoinPath(serviceAddress, s.ApiRoot)
	if err != nil {
		return err
	}
	s.ApiAddress = apiAddress

	return nil
}
