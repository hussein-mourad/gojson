package main

import (
	"os/exec"
	"path"
	"testing"
)

var testDir = "./tests"

func TestStep1Valid(t *testing.T) {
	t.Parallel()
	cmdValid := exec.Command("go", "run", "main.go", path.Join(testDir, "step1", "valid.json"))
	output, err := cmdValid.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	outputStr := string(output)
	expectedOutput := "{}"
	if outputStr != expectedOutput {
		t.Errorf("Expected %s but got %s", expectedOutput, outputStr)
	}
}
