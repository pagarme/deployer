package pipeline

type BuildStep struct {
	Config map[string]interface{}
}

func (s *BuildStep) Execute(p Context) error {
	return nil
}
