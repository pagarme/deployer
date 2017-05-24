package build

import (
	"errors"

	"github.com/pagarme/deployer/pipeline"
)

type BuildStep struct {
	Config map[string]interface{}
}

func (s *BuildStep) Execute(ctx pipeline.Context) error {
	typ, ok := s.Config["type"].(string)

	if !ok {
		return errors.New("missing type key")
	}

	builder, err := New(typ, s.Config)

	if err != nil {
		return err
	}

	return builder.Build(ctx)
}
