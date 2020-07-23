package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Pipeline     []map[string]map[string]interface{} `yml:"pipeline"`
	Deploy       map[string]interface{}              `yml:"deploy"`
	Environments map[string]interface{}              `yml:"environments"`
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

	m, ok := e.(map[interface{}]interface{})

	if !ok {
		return map[string]interface{}{}
	}

	real := map[string]interface{}{}

	for k, v := range m {
		real[fmt.Sprintf("%v", k)] = v
	}

	return real
}
