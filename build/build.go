package builder

import "errors"

type Factory func(config map[string]interface{}) (Builder, error)

type Builder interface {
}

var factories = map[string]Factory{}

func Register(name string, factory Factory) {
	factories[name] = factory
}

func New(name string, config map[string]interface{}) (Builder, error) {
	factory, ok := factories[name]

	if !ok {
		return nil, errors.New("invalid builder")
	}

	return factory(config)
}
