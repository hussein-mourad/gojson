package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestStepOne(t *testing.T) {
	runStepTests(t, "testdata/step1/*")
}

func TestStepTwo(t *testing.T) {
	runStepTests(t, "testdata/step2/*")
}

func TestStepThree(t *testing.T) {
	runStepTests(t, "testdata/step3/*")
}

func TestStepFour(t *testing.T) {
	runStepTests(t, "testdata/step4/*")
}

func TestStepFinal(t *testing.T) {
	runStepTests(t, "testdata/final/*")
}

func runStepTests(t *testing.T, filePathGlob string) {
	// t.Parallel()
	// Read all files in the testdata directory
	files, err := filepath.Glob(filePathGlob)
	if err != nil {
		t.Fatalf("failed to read testdata directory: %v", err)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			expectedExitCode := determineExpectedExitCode(file)
			fmt.Printf("expectedExitCode: %v\n", expectedExitCode)

			cmd := exec.Command("go", "run", "main.go", file)
			err := cmd.Run()

			exitCode := cmd.ProcessState.ExitCode()
			fmt.Printf("exitCode: %v\n", exitCode)
			if exitCode != expectedExitCode {
				t.Errorf("for %s, got exit code %d, want %d", file, exitCode, expectedExitCode)
			}

			if err != nil && exitCode == 0 {
				t.Fatalf("command failed unexpectedly: %v", err)
			}
		})
	}
}

// Helper function to determine expected exit code based on filename
func determineExpectedExitCode(filename string) int {
	lowerFilename := strings.ToLower(filename)
	if strings.Contains(lowerFilename, "fail") || strings.Contains(lowerFilename, "invalid") {
		return 1 // Any non-zero exit code for fail or invalid
	}
	return 0 // Exit code 0 for pass or valid
}
