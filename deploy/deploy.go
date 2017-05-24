package deployer

import "errors"

type Factory func(config map[string]interface{}) (Deployer, error)

type Deployer interface {
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
