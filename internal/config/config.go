package config

import (
	"fmt"
	"reflect"
)

type ResourceConfig interface {
	DownloadClient | RootFolder | ProwlarrApplication
	GetID() int
	GetName() string
	GetCreateData() map[string]interface{}
}

type RootFolder struct {
	ID   int    `json:"id"`
	Path string `mapstructure:"path" json:"path"`
}

func (rf RootFolder) GetID() int {
	return rf.ID
}

func (rf RootFolder) GetName() string {
	return rf.Path
}

func (rf RootFolder) GetCreateData() map[string]interface{} {
	return map[string]interface{}{"path": rf.Path}
}

type DownloadClient struct {
	ID      int             `json:"id"`
	Name    string          `mapstructure:"name"`
	AppType string          `mapstructure:"appType"`
	Fields  []ResourceField `json:"fields" mapstructure:"fields"`
}

func (dc DownloadClient) GetID() int {
	return dc.ID
}

func (dc DownloadClient) GetName() string {
	return dc.Name
}

func (dc DownloadClient) GetCreateData() map[string]interface{} {
	return map[string]interface{}{"enable": true, "name": dc.Name, "implementation": dc.AppType, "configContract": fmt.Sprintf("%sSettings", dc.AppType), "fields": dc.Fields}
}

type ProwlarrApplication struct {
	ID      int             `json:"id"`
	Name    string          `mapstructure:"name"`
	AppType string          `mapstructure:"appType"`
	URL     string          `mapstructure:"url"`
	Fields  []ResourceField `json:"fields" mapstructure:"fields"`
}

func (pa ProwlarrApplication) GetID() int {
	return pa.ID
}

func (pa ProwlarrApplication) GetName() string {
	return pa.Name
}

func (pa ProwlarrApplication) GetCreateData() map[string]interface{} {
	return map[string]interface{}{"syncLevel": "fullSync", "name": pa.Name, "implementation": pa.AppType, "configContract": fmt.Sprintf("%sSettings", pa.AppType), "fields": pa.Fields}
}

type ResourceField struct {
	Name  string      `json:"name" mapstructure:"name"`
	Value interface{} `json:"value" mapstructure:"value"`
}

type ServiceConfig struct {
	RootFolders     []RootFolder          `mapstructure:"rootfolder"`
	DownloadClients []DownloadClient      `mapstructure:"downloadClient"`
	Applications    []ProwlarrApplication `mapstructure:"applications"`
}

type Service struct {
	Name           string
	ApiAddress     string
	ServiceAddress string
	Address        string        `mapstructure:"address"`
	Port           uint16        `mapstructure:"port"`
	ApiRoot        string        `json:"apiRoot" yaml:"apiRoot"`
	ApiKey         string        `json:"apiKey" yaml:"apiKey"`
	Config         ServiceConfig `mapstructure:"config"`
}

func (s *Service) IsEmpty() bool {
	return reflect.DeepEqual(*s, Service{})
}

type Config struct {
	Services []Service `mapstructure:"services"`
}
