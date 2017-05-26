package scm

import (
	"errors"

	"github.com/pagarme/deployer/pipeline"
)

type Factory func(config map[string]interface{}) (Scm, error)

type Scm interface {
	Fetch(ctx pipeline.Context, workdir, ref string) error
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
