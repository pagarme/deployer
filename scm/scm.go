package scm

import "errors"

type Factory func(config map[string]interface{}) (Scm, error)

type Scm interface {
	Fetch(workdir, ref string) error
}

var factories = map[string]Factory{}

func Register(name string, factory Factory) {
	factories[name] = factory
}

func New(name string, config map[string]interface{}) (Scm, error) {
	factory, ok := factories[name]

	if !ok {
		return nil, errors.New("invalid scm")
	}

	return factory(config)
}
