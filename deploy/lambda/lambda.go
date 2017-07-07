package lambda

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mholt/archiver"
	"github.com/mitchellh/mapstructure"

	"github.com/pagarme/deployer/deploy"
	"github.com/pagarme/deployer/pipeline"
	"github.com/pagarme/deployer/scm"
)

type Options struct {
	S3Bucket    string   `mapstructure:"s3_bucket"`
	Region      string   `mapstructure:"region"`
	ProjectName string   `mapstructure:"project_name"`
	Package     []string `mapstructure:"package"`
}

type Environment struct {
	Name      string   `mapstructure:"name"`
	Functions []string `mapstructure:"functions"`
}

type Lambda struct {
	Environment *Environment
	Options     *Options
}

func (l *Lambda) zipPackage(path, filename string) (string, error) {
	var packages []string
	packagePath := filepath.Join(path, filename)

	for _, p := range l.Options.Package {
		packages = append(packages, filepath.Join(path, p))
	}

	if err := archiver.Zip.Make(packagePath, packages); err != nil {
		return "", err
	}

	return packagePath, nil
}

func (l *Lambda) uploadZipToS3(filename, S3Key string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(l.Options.Region)},

		SharedConfigState: session.SharedConfigEnable,
	}))

	uploader := s3manager.NewUploader(sess)
	file, err := os.Open(filename)

	if err != nil {
		return errors.New("could not open zip file")
	}

	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(l.Options.S3Bucket),
		Key:    aws.String(S3Key),
		Body:   file,
	})

	if err != nil {
		return err
	}

	return nil
}

func (l *Lambda) updateFunctionsCode(functions []string, S3Key string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(l.Options.Region)},

		SharedConfigState: session.SharedConfigEnable,
	}))

	lambdaService := lambda.New(sess)

	var wg sync.WaitGroup
	wg.Add(len(functions))

	doneChannel := make(chan bool)
	errChannel := make(chan error)

	for _, function := range functions {
		lambdaInput := &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(function),
			S3Bucket:     aws.String(l.Options.S3Bucket),
			S3Key:        aws.String(S3Key),
		}

		go func(input *lambda.UpdateFunctionCodeInput) {
			defer wg.Done()

			if _, err := lambdaService.UpdateFunctionCode(input); err != nil {
				errChannel <- err
			}
		}(lambdaInput)
	}

	go func() {
		wg.Wait()
		doneChannel <- true
	}()

	select {
	case <-doneChannel:
		return nil
	case err := <-errChannel:
		return err
	}
}

func (l *Lambda) Deploy(context pipeline.Context) error {
	environmentContext, ok := context["Environment"].(map[string]interface{})

	if !ok {
		return fmt.Errorf("could not get %s from %v map", "Environment", context)
	}

	environment := &Environment{}

	if err := mapstructure.Decode(environmentContext, environment); err != nil {
		return err
	}

	l.Environment = environment

	hash := "latest"

	if commitable, ok := context["Scm"].(scm.Commitable); ok {
		hash = commitable.CommitHash()

		if len(hash) > 7 {
			hash = hash[:8]
		}
	}

	zipName := "lambda.zip"
	S3Key := fmt.Sprintf("%s/%s/%s", l.Environment.Name, hash, zipName)

	ScmPath, ok := context["ScmPath"].(string)

	if !ok {
		return fmt.Errorf("could not get %s from %v map", "ScmPath", context)
	}

	filename, err := l.zipPackage(ScmPath, zipName)

	if err != nil {
		return err
	}

	if err := l.uploadZipToS3(filename, S3Key); err != nil {
		return err
	}

	if err := l.updateFunctionsCode(l.Environment.Functions, S3Key); err != nil {
		return err
	}

	return nil
}

func init() {
	deploy.Register("lambda", func(config map[string]interface{}) (deploy.Deployer, error) {
		options := &Options{}

		if err := mapstructure.Decode(config, options); err != nil {
			return nil, err
		}

		return &Lambda{Options: options}, nil
	})
}
