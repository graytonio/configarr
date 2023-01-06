package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/graytonio/configarr/internal/log"
	"github.com/sirupsen/logrus"
)

var client *http.Client

func init() {
	client = &http.Client{}
}

func GetResource(address string, apiKey string, resource string) ([]byte, error) {
	resourcePath, err := url.JoinPath(address, resource)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", resourcePath, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", apiKey)
	log.Logger.WithFields(logrus.Fields{
		"url":    resourcePath,
		"apiKey": apiKey,
	}).Debug("GET")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("could not GET %s: %s", resourcePath, resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, err
}

func PostResource(address string, apiKey string, resource string, data map[string]interface{}) error {
	resourcePath, err := url.JoinPath(address, resource)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", resourcePath, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	log.Logger.WithFields(logrus.Fields{
		"url":    resourcePath,
		"apiKey": apiKey,
		"body":   string(jsonData),
	}).Debug("POST")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("could not POST %s: %s", resourcePath, resp.Status)
	}

	return nil
}

func DeleteResource(address string, apiKey string, resource string, id string) error {
	resourcePath, err := url.JoinPath(address, resource, id)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", resourcePath, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", apiKey)
	log.Logger.WithFields(logrus.Fields{
		"url":    resourcePath,
		"apiKey": apiKey,
	}).Debug("DELETE")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("could not DELETE %s: %s", resourcePath, resp.Status)
	}

	return nil
}
