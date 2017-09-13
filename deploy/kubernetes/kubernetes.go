package kubernetes

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/mitchellh/mapstructure"

	"github.com/pagarme/deployer/deploy"
	"github.com/pagarme/deployer/pipeline"
)

type Options struct {
	DeploymentDir  string `mapstructure:"deployment_dir"`
	DeploymentFile string `mapstructure:"deployment_file"`
	Namespace      string `mapstructure:"namespace"`
}

type Kubernetes struct {
	Options *Options
}

func init() {
	deploy.Register("kubernetes", func(config map[string]interface{}) (deploy.Deployer, error) {
		options := &Options{}

		if err := mapstructure.Decode(config, options); err != nil {
			return nil, err
		}

		return &Kubernetes{Options: options}, nil
	})
}

func (n *Kubernetes) Deploy(ctx pipeline.Context) error {
	workdir, err := ioutil.TempDir("/tmp", "deployer")

	if err != nil {
		return err
	}

	if len(n.Options.DeploymentFile) != 0 {
		err := n.compileFile(ctx, n.Options.DeploymentFile, path.Join(workdir, "a.yaml"))

		if err != nil {
			return err
		}
	} else {
		err := n.compileDir(ctx, n.Options.DeploymentDir, workdir)

		if err != nil {
			return err
		}
	}

	err = n.apply(ctx, workdir)

	if err != nil {
		return err
	}

	return nil
}

func (n *Kubernetes) compileDir(ctx pipeline.Context, in, out string) error {
	files, err := ioutil.ReadDir(in)

	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}

		inFile := path.Join(in, f.Name())
		outFile := path.Join(out, f.Name())

		err = n.compileFile(ctx, inFile, outFile)

		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Kubernetes) compileFile(ctx pipeline.Context, in, out string) error {
	w, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	t, err := template.New(in).ParseFiles(in)

	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(w, path.Base(in), ctx)

	if err != nil {
		return err
	}

	return w.Close()
}

func (n *Kubernetes) apply(ctx pipeline.Context, in string) error {
	cmd := exec.Command("kubectl", "-n", n.Options.Namespace, "apply", "-f", in)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
