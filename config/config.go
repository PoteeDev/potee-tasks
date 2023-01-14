package config

import (
	"console/models"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Groups         []string           `yaml:"groups"`
	ServicesSubnet string             `yaml:"services_subnet"`
	ClinetsSubnet  string             `yaml:"clients_subnet"`
	IpCount        int                `yaml:"ip_count"`
	Categories     []string           `yaml:"categories"`
	Tasks          []models.Challenge `yaml:"tasks"`
}

func ReadConfig(filename string) *Config {
	yfile, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err = yaml.Unmarshal(yfile, &config); err != nil {

		log.Fatal(err)
	}
	return &config
}
