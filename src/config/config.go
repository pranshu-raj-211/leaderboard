package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Redis struct {
		Address  string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
}

var AppConfig *Config

func LoadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	AppConfig = &cfg
	return nil
}
