package scm

import (
	"errors"
	"fmt"

	"github.com/pagarme/deployer/pipeline"
)

type ScmStep struct {
	Config map[string]interface{}
	Ref    string
}

func (s *ScmStep) Execute(ctx pipeline.Context) error {
	typ, ok := s.Config["type"].(string)

	if !ok {
		return errors.New("missing type key")
	}

	scmInstance, err := New(typ, s.Config)

	if err != nil {
		return err
	}

	// workdir, err := ioutil.TempDir("/tmp", "deployer")
	//
	// if err != nil {
	// 	return err
	// }

	fmt.Println("Remember to create /tmp/superbowleto-deployer")
	workdir := "/tmp/superbowleto-deployer"

	ctx["ScmPath"] = workdir

	return scmInstance.Fetch(ctx, workdir, s.Ref)
}
