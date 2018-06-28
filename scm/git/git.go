package git

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"

	"github.com/mitchellh/mapstructure"
	"github.com/pagarme/deployer/pipeline"
	"github.com/pagarme/deployer/scm"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Options struct {
	Repository string `mapstructure:"repository"`
	Ref        string `mapstructure:"ref"`
}

type Git struct {
	Options *Options
}

type Context struct {
	Repo   *git.Repository
	Commit plumbing.Hash
}

func init() {
	scm.Register("git", func(config map[string]interface{}) (scm.Scm, error) {
		options := &Options{}

		if err := mapstructure.Decode(config, options); err != nil {
			return nil, err
		}

		return &Git{Options: options}, nil
	})
}

func (g *Git) Fetch(ctx pipeline.Context, workdir, refName string) error {
	var hash plumbing.Hash

	repo, err := git.PlainClone(workdir, false, &git.CloneOptions{
		RemoteName: "origin",
		URL:        g.Options.Repository,
		Progress:   os.Stdout,
	})

	if err != nil {
		return err
	}

	isSha, err := regexp.MatchString("\\b[0-9a-fA-F]{5,40}\\b", refName)

	if err != nil {
		return err
	}

	ref, err := repo.Reference(plumbing.ReferenceName("refs/tags/"+refName), true)

	if err != nil {
		return err
	}

	if ref != nil {
		hash = ref.Hash()
	} else if isSha {
		hash = plumbing.NewHash(refName)
	} else {
		return fmt.Errorf("invalid ref")
	}

	wt, err := repo.Worktree()

	if err != nil {
		return err
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Hash: hash,
	})

	if err != nil {
		return err
	}

	ctx["Scm"] = &Context{
		Repo:   repo,
		Commit: hash,
	}

	return nil
}

func (c *Context) CommitHash() string {
	return hex.EncodeToString(c.Commit[:])
}
