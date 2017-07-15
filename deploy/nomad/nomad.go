package nomad

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/mitchellh/mapstructure"

	"github.com/pagarme/deployer/deploy"
	"github.com/pagarme/deployer/pipeline"
)

const (
	EnvNomadAddress = "NOMAD_ADDR"
	EnvNomadRegion  = "NOMAD_REGION"
)

type Options struct {
	JobFile  interface{}   `mapstructure:"job_file"`
	JobFiles []interface{} `mapstructure:"job_files"`
}

type Jobfile struct {
	File           string `mapstructure:"file"`
	WaitAllocation bool   `mapstructure:"wait_allocation"`
}

type Nomad struct {
	Options *Options

	JobFiles []Jobfile
}

func init() {
	deploy.Register("nomad", func(config map[string]interface{}) (deploy.Deployer, error) {
		options := &Options{}

		if err := mapstructure.Decode(config, options); err != nil {
			return nil, err
		}

		jobFiles := append(options.JobFiles, options.JobFile)
		jobs := make([]Jobfile, len(jobFiles))

		for i, v := range jobFiles {
			var job Jobfile

			if v == nil {
				continue
			}

			switch v := v.(type) {
			case string:
				job = Jobfile{File: v, WaitAllocation: true}
			case map[interface{}]interface{}:
				if err := mapstructure.Decode(v, &job); err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("invalid configuration format %v", v)
			}

			jobs[i] = job
		}

		if options.JobFile == "" && len(options.JobFiles) == 0 {
			options.JobFile = "deploy.nomad"
		}

		return &Nomad{Options: options, JobFiles: jobs}, nil
	})
}

func (n *Nomad) Deploy(ctx pipeline.Context) error {
	for _, v := range n.JobFiles {
		jobSource, err := n.compileJob(ctx, v)

		if err != nil {
			return err
		}

		job, err := jobspec.Parse(jobSource)

		if err != nil {
			return err
		}

		client, err := n.client()

		if err != nil {
			return err
		}

		evalID, _, err := client.Jobs().Register(job, nil)

		if err != nil {
			return err
		}

		for v.WaitAllocation {
			eval, _, err := client.Evaluations().Info(evalID, nil)

			if err != nil {
				return err
			}

			switch eval.Status {
			case structs.EvalStatusFailed, structs.EvalStatusCancelled:
				return fmt.Errorf("allocation failed")

			case structs.EvalStatusComplete:
				return n.watchJob(client, eval)

			default:
				time.Sleep(time.Second)
			}
		}
	}

	return nil
}

func (n *Nomad) watchJob(client *api.Client, eval *api.Evaluation) error {
OUTER:
	for true {
		allocs, _, err := client.Jobs().Allocations(eval.JobID, false, nil)

		if err != nil {
			return err
		}

		evalAllocs := []*api.AllocationListStub{}

		for _, a := range allocs {
			if a.EvalID == eval.ID {
				evalAllocs = append(evalAllocs, a)
			}
		}

		if len(evalAllocs) == 0 {
			return fmt.Errorf("no allocation was made")
		}

		for _, alloc := range evalAllocs {
			switch alloc.ClientStatus {
			case structs.AllocClientStatusFailed, structs.AllocClientStatusLost, structs.AllocClientStatusComplete:
				return fmt.Errorf("allocation %s failed", alloc.ID)

			case structs.AllocClientStatusRunning:
				break OUTER

			default:
				time.Sleep(time.Second)
				continue OUTER
			}
		}
	}

	return nil
}

func (n *Nomad) compileJob(ctx pipeline.Context, job Jobfile) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)

	t, err := template.ParseFiles(job.File)

	if err != nil {
		return nil, err
	}

	err = t.ExecuteTemplate(buf, job.File, ctx)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(buf.Bytes()), nil
}

func (n *Nomad) client() (*api.Client, error) {
	config := api.DefaultConfig()

	if v := os.Getenv(EnvNomadAddress); v != "" {
		config.Address = v
	}

	if v := os.Getenv(EnvNomadRegion); v != "" {
		config.Region = v
	}

	return api.NewClient(config)
}
