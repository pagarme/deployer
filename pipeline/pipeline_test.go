package pipeline

import (
	"errors"
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	result := Create()
	expectedResult := &Pipeline{
		Context: make(Context),
		steps:   make([]Step, 0),
	}

	if reflect.ValueOf(result) == reflect.ValueOf(expectedResult) {
		t.Error("Expected:", expectedResult, "\nGot:", result)
	}
	if reflect.ValueOf(result.Context) == reflect.ValueOf(expectedResult.Context) {
		t.Error("Expected:", expectedResult, "\nGot:", result)
	}
	if reflect.ValueOf(result.steps) == reflect.ValueOf(expectedResult.steps) {
		t.Error("Expected:", expectedResult, "\nGot:", result)
	}
}

func TestAdd(t *testing.T) {
	pipeline := Create()
	var step Step
	pipeline.Add(step)

	expectedPipeline := Create()
	var expectedStep Step
	expectedPipeline.steps = append(expectedPipeline.steps, expectedStep)

	if !reflect.DeepEqual(expectedPipeline.steps, pipeline.steps) {
		t.Error("Expected:", expectedPipeline.steps, "\nGot:", pipeline.steps)
	}
}

type TestStep struct {
}

func (t TestStep) Execute(p Context) error {
	return nil
}

type ErrorStep struct {
}

func (t ErrorStep) Execute(p Context) error {
	return errors.New("Step Error!")
}

func TestExecute(t *testing.T) {
	pipeline := Create()
	var step TestStep
	pipeline.Add(step)
	result := pipeline.Execute()

	if result != nil {
		t.Error("Expected:", nil, "\nGot:", result)
	}

	pipelineError := Create()
	var stepError ErrorStep
	pipelineError.Add(stepError)
	resultError := pipelineError.Execute()

	if resultError == nil {
		t.Error("Expected:", nil, "\nGot:", result)
	}
}
