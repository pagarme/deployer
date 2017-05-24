package pipeline

type Step interface {
	Execute(p Context) error
}
