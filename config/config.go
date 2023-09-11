package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var config Config

type Config struct {
	Identity string `yaml:"identity"`
	Cards    []Card `yaml:"cards"`

	Goodreads Goodreads `yaml:"goodreads"`

	path string
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

type Goodreads struct {
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token"`
	ExpiresAt    int64  `yaml:"expires_at"`
	UserId       string `yaml:"user_id"`
}

func Get() *Config {
	return &config
}

func Load(filename string) error {
	newConfig := Config{
		path: filename,
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&newConfig)
	if err != nil {
		return err
	}

	config = newConfig

	return nil
}

func (c *Config) Save() error {
	f, err := os.OpenFile(c.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	return yaml.NewEncoder(f).Encode(config)
}
