package pipeline

type DeployStep struct {
	Config map[string]interface{}
}

func (s *DeployStep) Execute(p Context) error {
	return nil
}
