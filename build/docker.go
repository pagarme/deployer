package build

type Docker interface {
	DockerImage() string
}
