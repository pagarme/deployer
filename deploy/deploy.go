package deploy

import (
	"errors"

	"github.com/pagarme/deployer/pipeline"
)

type Factory func(config map[string]interface{}) (Deployer, error)

type Deployer interface {
	Deploy(ctx pipeline.Context) error
}

var factories = map[string]Factory{}

func Register(name string, factory Factory) {
	factories[name] = factory
}

func New(name string, config map[string]interface{}) (Deployer, error) {
	factory, ok := factories[name]

	if !ok {
		return nil, errors.New("invalid deployer")
	}

	return factory(config)
}
