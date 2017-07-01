# deployer

:pager: A tool for fetching, building, pushing and deploying applications.

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
  --ref     Source Code Management hash to be fetched (default: master)
  --env     Environment to be used (default: main)
```

## Configuration File

To deploy an application you must specify a yml configuration file (e.g. `deployer.yml`), that consists in `steps` and `environment` configuration.
A typical configuration file has the following structure:

```yml
scm:
  type: <type>

build:
  type: <type>

deploy:
  type: <type>

environments:
  sandbox:
    name: sandbox
  live:
    name: live
```

**Note:** The order the steps appear in the configuration file, does not determine the order they will be executed. Check [Steps](#steps) for more information.

## Steps

The `deployer` deploys application executing different steps, one after the other. The order the steps are executed is: `scm -> build -> deploy`

## SCM (Source Code Management)

Available types are: `git`.

### Git

The Git SCM Step clones a git repository with an specified hash and adds it to a temporary folder.

**Options:**

| Options | Description |
| --- | --- |
| `repository` | The git repository to be cloned |
| `ref` | The ref of the repository to be used. Can be a branch or the sha of a commit. Defaults to `master` or the ref passed via flag `--ref` |

**Example:**

```yml
scm:
  type: git
  repository: https://github.com/pagarme/deployer
  ref: master
```

## Build

Available types are: `rocker`.

### Rocker

The Rocker Build step builds and pushes a specified Rockerfile. The following variables are injected when building the Rockerfile:

  1. `RepositoryPath`: The location of the cloned repository from the SCM step
  2. `ImageSHA`: The sha of the ref from the SCM step
  3. `ImageRepository`: The repository where the Rocker image will be pushed

For more information about Rocker, check the official documentation (https://github.com/grammarly/rocker)

**Options:**

| Options | Description |
| --- | --- |
| `build_directory` | The location of the Rockerfile. Defaults to current directory |
| `image_repository` | The repository where the Rocker image will be pushed |

**Example:**

```yml
build:
  type: rocker
  build_directory: path/to
  image_repository: xxxxxx.dkr.ecr.us-east-1.amazonaws.com/deployer # Pushing the image to an AWS ECR
```

## Deploy

Available types are: `lambda` and `nomad`.

### Lambda

The Lambda Deploy step updates and deploys existing AWS Lambda functions. It creates a zip file, uploads it to a S3 Bucket and updates the Lambda functions code.

There's a convention that must be followed to name the Lambda functions. The functions names must have the following shape: `env-projectName-functionName`. So a function called `hello-world` from the `pagarme` project on the `live` environment must be called `live-pagarme-hello-world.

**Note:** This step does not create the AWS Lambda functions, but instead updates existing functions.

**Options:**

| Options | Description |
| --- | --- |
| `project_name` | The name of the project. This will prefix the functions names. |
| `region` | The AWS region the functions are deployed. |
| `s3_bucket` | The name of the AWS S3 Bucket that will be used to upload the functions zipped source code |
| `package` | The source code of the functions. A list of folders and files to include in the zip. |
| `functions` | The names of the functions (without the `env` and `project_name` prefix. |

**Example:**

```yml
deploy:
  type: lambda
  project_name: pagarme
  region: us-east-1
  s3_bucket: pagarme-deploy-s3-bucket
  package: # This will create a zip containing the `dist` and the `node_modules` folders
    - dist
    - node_modules
  functions: # The function names to be prefixed in the format `env-project-function`
    - hello  # The full name of the function is `live-pagarme-hello`
    - cowsay # The full name of the function is `live-pagarme-cowsay`
```

##
<p align="center">  
    <a href="https://github.com/pagarme" style="text-decoration:none; margin-right:2rem;">
    <img src="https://cdn.rawgit.com/pagarme/brand/9ec30d3d4a6dd8b799bca1c25f60fb123ad66d5b/logo-circle.svg" width="110px" height="110px" />
  </a>
</p>
