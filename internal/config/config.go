package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type changeType string

const (
	Add    changeType = "add"
	Remove changeType = "remove"
	Modify changeType = "modify"
)

type HttpVerb string

const (
	GET    HttpVerb = "GET"
	POST   HttpVerb = "POST"
	PUT    HttpVerb = "PUT"
	PATCH  HttpVerb = "PATCH"
	DELETE HttpVerb = "DELETE"
)

type Change struct {
	Type  changeType  `yaml:"type"`
	Field string      `yaml:"field"`
	Value interface{} `yaml:"value,omitempty"`
}

type EndpointConfig struct {
	Path        string   `yaml:"path"`
	Verb        HttpVerb `yaml:"verb"`
	StatusCodes []int    `yaml:"statusCodes"`
	Changes     []Change `yaml:"changes"`
}

type ServiceConfig struct {
	Host      string           `yaml:"host"`
	Endpoints []EndpointConfig `yaml:"endpoints"`
}

type MactConfig struct {
	Services []ServiceConfig `yaml:"services"`
}

func ReadConfig(configFile string) MactConfig {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error readig config file %v: %v ", configFile, err)
	}
	var config MactConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return config
}
