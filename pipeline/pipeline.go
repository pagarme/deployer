package pipeline

type Pipeline struct {
	Context Context
	steps   []Step
}

func Create() *Pipeline {
	return &Pipeline{
		Context: make(Context),
		steps:   make([]Step, 0),
	}
}

func (p *Pipeline) Add(step Step) {
	p.steps = append(p.steps, step)
}

func (p *Pipeline) Execute() error {
	for _, s := range p.steps {
		err := s.Execute(p.Context)

		if err != nil {
			return err
		}
	}

	return nil
}
