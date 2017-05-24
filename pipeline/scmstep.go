package pipeline

import (
	"errors"
	"io/ioutil"

	"github.com/pagarme/deployer/scm"
)

type ScmStep struct {
	Config map[string]interface{}
	Ref    string
}

func (s *ScmStep) Execute(ctx Context) error {
	typ, ok := s.Config["type"].(string)

	if !ok {
		return errors.New("missing type key")
	}

	scmInstance, err := scm.New(typ, s.Config)

	if err != nil {
		return err
	}

	workdir, err := ioutil.TempDir("", "deployer")

	if err != nil {
		return err
	}

	ctx["ScmWorkingDirectory"] = workdir

	print(workdir)
	print("\n")

	return scmInstance.Fetch(workdir, s.Ref)
}
