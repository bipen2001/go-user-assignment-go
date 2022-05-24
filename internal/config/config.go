package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Port       int    `yaml:"port"`
	Host       string `yaml:"host"`
	CorsOrigin string `yaml:"corsOrigin"`
}
type Db struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
}
type Auth struct {
	JwtKey string `yaml:"jwtkey"`
}

type Pagination struct {
	Limit       int `yaml:"Limit"`
	DefaultPage int `yaml:"defaultpage"`
}
type Config struct {
	Server `yaml:"server"`

	Database Db `yaml:"database"`

	Auth `yaml:"auth"`

	Pagination `yaml:"pagination"`
}

func Load() (*Config, error) {

	appConfig := &Config{}

	configFile := "local.yaml"

	if _, err := os.Stat("../../" + configFile); err != nil {

		log.Printf("could not find local.yaml in directory: %s", err)

		configFile = "config.yaml"
	}

	bytes, err := ioutil.ReadFile("../../" + configFile)
	if err != nil {

		log.Printf("error reading config file: %s", err)

		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &appConfig); err != nil {

		log.Printf("error unmarshalling config: %s", err)

		return nil, err
	}

	return appConfig, nil

}
