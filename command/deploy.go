package command

import (
	"flag"
	"fmt"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/pagarme/deployer/build"
	"github.com/pagarme/deployer/config"
	"github.com/pagarme/deployer/deploy"
	"github.com/pagarme/deployer/logger"
	"github.com/pagarme/deployer/pipeline"
	"github.com/pagarme/deployer/scm"
	uuid "github.com/satori/go.uuid"
)

type DeployCommand struct {
}

func (c *DeployCommand) Help() string {
	helpText := `
Usage:
	deployer command [options] <path>

Available Commands:
  deploy    Deploy an application using a configuration file

Options:
  --ref     Source Code Management hash to be fetched (default: master)
  --env     Environment to be used (default: main)
`

	return strings.TrimSpace(helpText)
}

func (c *DeployCommand) Synopsis() string {
	return "Executes a deploy"
}

func (c *DeployCommand) Run(args []string) int {
	log := &logger.DynamoLogger{}
	log.Init()

	executionID := uuid.NewV4().String()
	curUser, _ := user.Current()

	commandLog := &logger.CommandLog{
		Username:     curUser.Username,
		Timestamp:    strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		Command:      "deploy",
		Args:         args,
		Status:       "started",
		StatusReason: "none",
		ExecutionID:  executionID,
	}

	log.LogCommand(*commandLog)

	env := ""
	ref := "master"

	flags := flag.NewFlagSet("deployer deploy", flag.ContinueOnError)
	flags.StringVar(&env, "env", "main", "Deployment environment")
	flags.StringVar(&ref, "ref", "master", "Scm reference to deploy")

	if err := flags.Parse(args); err != nil {
		fmt.Println(c.Help())
		return 1
	}

	args = flags.Args()

	if len(args) != 1 {
		fmt.Println(c.Help())
		return 1
	}

	cfg, err := config.ReadConfig(args[0])
	if err != nil {
		fmt.Printf("An error ocurred reading the configuration: %s\n", err)
		return 1
	}

	pipe := pipeline.Create()

	pipe.Context["Config"] = cfg
	pipe.Context["Environment"] = cfg.GetEnvironment(env)

	if len(cfg.Pipeline) == 0 {
		pipe.Add(&scm.ScmStep{
			Config: cfg.Scm,
			Ref:    ref,
		})

		pipe.Add(&build.BuildStep{Config: cfg.Build})

		pipe.Add(&deploy.DeployStep{Config: cfg.Deploy})
	} else {
		for _, v := range cfg.Pipeline {
			kind := ""

			for k := range v {
				if k != "scm" && k != "build" && k != "deploy" {
					continue
				}

				kind = k
				break
			}

			switch kind {
			case "scm":
				pipe.Add(&scm.ScmStep{
					Config: v[kind],
					Ref:    ref,
				})

			case "build":
				pipe.Add(&build.BuildStep{Config: v[kind]})

			case "deploy":
				pipe.Add(&deploy.DeployStep{Config: v[kind]})

			default:
				fmt.Printf("Invalid pipeline step\n")
				return 1
			}
		}
	}

	if err := pipe.Execute(); err != nil {
		fmt.Printf("An error ocurred during the deployment: %s\n", err)
		commandLog.Status = "failed"
		commandLog.StatusReason = err.Error()
		commandLog.Timestamp = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		log.LogCommand(*commandLog)

		return 1
	}

	commandLog.Status = "finished"
	commandLog.Timestamp = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	log.LogCommand(*commandLog)

	return 0
}
