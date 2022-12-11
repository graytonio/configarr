package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

type Service struct {
	Address string
	Port    uint16
	ApiRoot string `json:"apiRoot" yaml:"apiRoot"`
	ApiKey  string `json:"apiKey" yaml:"apiKey"`
}

type ServiceInitializeResponse struct {
	ApiRoot      string `json:"apiRoot" yaml:"apiRoot"`
	ApiKey       string `json:"apiKey" yaml:"apiKey"`
	Release      string `json:"release" yaml:"release"`
	Version      string `json:"version" yaml:"version"`
	InstanceName string `json:"instanceName" yaml:"instanceName"`
	Branch       string `json:"branch" yaml:"branch"`
	Analytics    bool   `json:"analytics" yaml:"analytics"`
	UrlBase      string `json:"urlBase" yaml:"urlBase"`
	IsProduction bool   `json:"isProduction" yaml:"isProduction"`
}

func NewService(address string, port uint16) (*Service, error) {

	var service = Service{
		Address: address,
		Port:    port,
	}

	err := service.initService()
	if err != nil {
		return nil, err
	}

	return &service, err
}

func (s *Service) initService() error {
	// Format url correctly
	serviceAddress := s.Address
	if !strings.HasPrefix(serviceAddress, "http") {
		serviceAddress = fmt.Sprintf("http://%s", serviceAddress)
	}
	serviceAddress = fmt.Sprintf("%s:%d", serviceAddress, s.Port)

	serviceAddress, err := url.JoinPath(serviceAddress, "initialize.js")
	if err != nil {
		return err
	}

	// Fetch initialization variables
	resp, err := http.Get(serviceAddress)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jsonData := strings.Trim(strings.Split(string(body), "=")[1], " ")

	err = yaml.Unmarshal([]byte(jsonData), &s)
	if err != nil {
		return err
	}

	return nil
}
