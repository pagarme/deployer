package rocker

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mitchellh/mapstructure"

	"github.com/pagarme/deployer/build"
	"github.com/pagarme/deployer/pipeline"
	"github.com/pagarme/deployer/scm"
)

type Options struct {
	BuildDirectory  string `mapatructure:"build_directory"`
	ImageRepository string `mapstructure:"image_repository"`
}

type Rocker struct {
	Options *Options
}

type Result struct {
	Image string
	Tag   string
}

func (r *Result) DockerImage() string {
	return fmt.Sprintf("%s:%s", r.Image, r.Tag)
}

func init() {
	build.Register("rocker", func(config map[string]interface{}) (build.Builder, error) {
		options := &Options{}

		if err := mapstructure.Decode(config, options); err != nil {
			return nil, err
		}

		if options.BuildDirectory == "" {
			wd, err := os.Getwd()

			if err != nil {
				return nil, err
			}

			options.BuildDirectory = wd
		}

		return &Rocker{Options: options}, nil
	})
}

func (r *Rocker) Build(ctx pipeline.Context) error {
	hash := "latest"

	if commitable, ok := ctx["Scm"].(scm.Commitable); ok {
		hash = commitable.CommitHash()

		if len(hash) > 7 {
			hash = hash[:8]
		}
	}

	args := []string{}

	args = append(args, "build")
	// args = append(args, "--no-cache")
	args = append(args, "--push")
	args = append(args, "--var")
	args = append(args, fmt.Sprintf("RepositoryPath=%s", ctx["ScmPath"]))
	args = append(args, "--var")
	args = append(args, fmt.Sprintf("ImageSHA=%s", hash))
	args = append(args, "--var")
	args = append(args, fmt.Sprintf("ImageRepository=%s", r.Options.ImageRepository))

	cmd := exec.Command("rocker", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = r.Options.BuildDirectory
	cmd.Env = os.Environ()

	err := cmd.Start()

	if err != nil {
		return err
	}

	err = cmd.Wait()

	if err != nil {
		return err
	}

	ctx["Build"] = &Result{
		Image: r.Options.ImageRepository,
		Tag:   hash,
	}

	return nil

}
