package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Scm          map[string]interface{} `yml:"scm"`
	Build        map[string]interface{} `yml:"build"`
	Deploy       map[string]interface{} `yml:"deploy"`
	Environments map[string]interface{} `yml:"environments"`
}

func ReadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) GetEnvironment(env string) map[string]interface{} {
	e, ok := c.Environments[env]

	if !ok {
		return map[string]interface{}{}
	}

	m, ok := e.(map[string]interface{})

	if !ok {
		return map[string]interface{}{}
	}

	return m
}
