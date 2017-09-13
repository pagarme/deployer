package scm

import (
	"errors"
	"io/ioutil"

	"github.com/pagarme/deployer/pipeline"
)

type ScmStep struct {
	Config map[string]interface{}
	Ref    string
}

func (s *ScmStep) Execute(ctx pipeline.Context) error {
	cfg, err := ctx.ResolveConfig(s.Config)

	if err != nil {
		return err
	}

	typ, ok := cfg["type"].(string)

	if !ok {
		return errors.New("missing type key")
	}

	scmInstance, err := New(typ, cfg)

	if err != nil {
		return err
	}

	workdir, err := ioutil.TempDir("/tmp", "deployer")

	if err != nil {
		return err
	}

	ctx["ScmPath"] = workdir

	return scmInstance.Fetch(ctx, workdir, s.Ref)
}
