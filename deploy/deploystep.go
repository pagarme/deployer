package deploy

import (
	"errors"

	"github.com/pagarme/deployer/pipeline"
)

type DeployStep struct {
	Config map[string]interface{}
}

func (s *DeployStep) Execute(ctx pipeline.Context) error {
	typ, ok := s.Config["type"].(string)

	if !ok {
		return errors.New("missing type key")
	}

	deployer, err := New(typ, s.Config)

	if err != nil {
		return err
	}

	return deployer.Deploy(ctx)
}
