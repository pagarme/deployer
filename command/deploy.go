package command

import (
	"flag"
	"fmt"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/pagarme/deployer/config"
	"github.com/pagarme/deployer/deploy"
	"github.com/pagarme/deployer/logger"
	"github.com/pagarme/deployer/pipeline"
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
  --env     Environment to be used (default: main)
  --img     Docker Image to be used
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
	img := ""

	flags := flag.NewFlagSet("deployer deploy", flag.ContinueOnError)
	flags.StringVar(&env, "env", "main", "Deployment environment")
	flags.StringVar(&img, "img", "main", "Docker image")

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
	pipe.Context["Image"] = img

	if len(cfg.Pipeline) == 0 {
		pipe.Add(&deploy.DeployStep{Config: cfg.Deploy})
	} else {
		for _, v := range cfg.Pipeline {
			kind := ""

			for k := range v {
				if k != "deploy" {
					continue
				}

				kind = k
				break
			}

			switch kind {
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
