package git

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/pagarme/deployer/scm"
)

type GitOptions struct {
	Repository string `mapstructure:"repository"`
	Ref        string `mapstructure:"ref"`
}

type Git struct {
	Options *GitOptions
}

func init() {
	scm.Register("git", func(config map[string]interface{}) (scm.Scm, error) {
		options := &GitOptions{}

		if err := mapstructure.Decode(config, options); err != nil {
			return nil, err
		}

		return &Git{Options: options}, nil
	})
}

func (g *Git) Fetch(workdir, ref string) error {
	_, err := git.PlainClone(workdir, false, &git.CloneOptions{
		URL:           g.Options.Repository,
		ReferenceName: plumbing.ReferenceName(ref),
		Progress:      os.Stdout,
	})

	return err
}
