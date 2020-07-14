# deployer

[![Build Status](https://travis-ci.org/pagarme/deployer.svg?branch=master)](https://travis-ci.org/pagarme/deployer)
[![Go Report Card](https://goreportcard.com/badge/github.com/pagarme/deployer)](https://goreportcard.com/report/github.com/pagarme/deployer)

:pager: A tool for deploying applications.

## Install

```sh
$ go get github.com/pagarme/deployer
```

## Usage

```
deployer command [options] <path>

Commands:
  deploy    Deploy an application using a configuration file

Options:
  --env     Environment to be used (default: main)
  --img     Docker Image to be used
```

## Configuration File

To deploy an application you must specify a yml configuration file (e.g. `deployer.yml`), that consists in `steps` and `environment` configuration.
A typical configuration file has the following structure:

```yml
deploy:
  type: <type>

environments:
  sandbox:
    name: sandbox
  live:
    name: live
```

**Note:** The order the steps appear in the configuration file, does not determine the order they will be executed. Check [Steps](#steps) for more information.

Also in order to send the logs to DynamoDB, the following environment 
variables must be set:
  - `DEPLOYER_DYNAMODB_TABLE`: DynamoDB's log table
  - `DEPLOYER_AWS_REGION`: AWS region 

## Deploy

Only `nomad`, example file:

### deploy.yml

```yml
---
deploy:
  type: nomad
  job_files:
    - application.nomad

environments:
  sandbox:
    name: sandbox
  live:
    name: live
```
### application.nomad

```json
job "application-{{ .Environment.name }}" {
  type        = "system"
  region      = "global"
  datacenters = ["dc1"]

  constraint {
    attribute = "${node.class}"
    value     = "application-{{ .Environment.name }}"
  }

  update {
    stagger      = "5s"
    max_parallel = 1
  }

  group "web" {
    task "nginx" {
      driver = "docker"

      config {
        image      = "{{ .Image }}"
        force_pull = true

        port_map {
          "http" = 80
        }

        dns_servers = ["${attr.unique.network.ip-address}"]
      }

      env {
        CONFIG_VERSION = "v1"
        CONSUL_ADDR    = "${attr.unique.network.ip-address}:8500"
      }

      service {
        name = "application-{{ .Environment.name }}"

        tags = ["application"]

        port = "http"

        check {
          name     = "http check"
          type     = "http"
          port     = "http"
          path     = "/_health_check"
          interval = "5s"
          timeout  = "2s"
        }

        check {
          name = "Status Information Health Check"
          type = "http"
          port = "http"
          path = "/_status"
          interval = "5s"
          timeout = "2s"
        }
      }

      resources {
        cpu    = "2400"
        memory = "3584"

        network {
          mbits = 100

          port "http" {
            static = 80
          }
        }
      }
    }
  }
}
```