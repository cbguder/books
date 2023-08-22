package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Identity string `yaml:"identity"`
	Cards    []Card `yaml:"cards"`
}

type Card struct {
	Id      string  `yaml:"id"`
	Name    string  `yaml:"name"`
	Library Library `yaml:"library"`
}

type Library struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

func ReadConfig(filename string) (*Config, error) {
	config := Config{}

	f, err := os.Open(filename)
	if err != nil {
		return &config, err
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&config)
	return &config, err
}

func WriteConfig(filename string, config *Config) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	return yaml.NewEncoder(f).Encode(config)
}
