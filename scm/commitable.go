package scm

type Commitable interface {
	CommitHash() string
}
