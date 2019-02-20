package rocker

import (
	"testing"
)

func TestDockerImage(t *testing.T) {
	results := []*Result{
		{},
		{"", ""},
		{"Image", ""},
		{"", "Tag"},
		{"Image", "Tag"},
	}

	var dockerImageStrs []string

	for _, result := range results {
		dockerImageStrs = append(dockerImageStrs, result.DockerImage())
	}

	var expectedResults = []string{
		":",
		":",
		"Image:",
		":Tag",
		"Image:Tag",
	}

	for i, dockerImageStr := range dockerImageStrs {
		if dockerImageStr != expectedResults[i] {
			t.Error("Expected:", expectedResults[i], "\nGot:", dockerImageStr)
		}
	}
}
