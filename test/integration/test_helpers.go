package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestContainer struct {
	Name        string
	Environment string
	DockerCmd   *exec.Cmd
}

type TestScenario struct {
	Name               string
	Environment        string
	ToolsToInstall     []string
	ExpectedBinaries   []string
	ExpectedConfigs    []string
	ValidationCommands []ValidationCmd
}

type ValidationCmd struct {
	Command     string
	Args        []string
	ExpectedOut string
	ShouldFail  bool
}

func SetupTestContainer(t *testing.T, environment string) *TestContainer {
	// Build the devenv binary
	buildCmd := exec.Command("go", "build", "-o", "bin/devenv")
	buildCmd.Dir = getProjectRoot()
	require.NoError(t, buildCmd.Run(), "Failed to build devenv binary")

	// Start the container
	// containerName := fmt.Sprintf("devenv-test-%s-%d", environment, time.Now().Unix())
	containerName := fmt.Sprintf("devenv-test-%s", environment)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Run docker compose
	composeCmd := exec.CommandContext(ctx, "docker compose",
		"-f", "test/integration/docker/docker-compose.yml",
		// "run", "--rm", "--name", containerName,
		"run", "--name", containerName,
		fmt.Sprintf("%s-test", environment))

	fmt.Printf("Starting container %s for environment %s\n", containerName, environment)
	fmt.Printf("Running command: %s\n", composeCmd.String())

	return &TestContainer{
		Name:        containerName,
		Environment: environment,
		DockerCmd:   composeCmd,
	}
}

func (tc *TestContainer) ExecuteCommand(t *testing.T, command string, args ...string) (string, error) {
	fullCmd := append([]string{"docker", "exec", tc.Name}, append([]string{command}, args...)...)

	cmd := exec.Command(fullCmd[0], fullCmd[1:]...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func (tc *TestContainer) ExecuteDevEnv(t *testing.T, args ...string) (string, error) {
	return tc.ExecuteCommand(t, "/usr/local/bin/devenv", args...)
}

func (tc *TestContainer) CheckBinaryExists(t *testing.T, binary string) bool {
	_, err := tc.ExecuteCommand(t, "which", binary)
	return err == nil
}

func (tc *TestContainer) CheckConfigExists(t *testing.T, configPath string) bool {
	_, err := tc.ExecuteCommand(t, "test", "-f", configPath)
	return err == nil
}

func (tc *TestContainer) ValidateInstallation(t *testing.T, scenario TestScenario) {
	// Check binaries exist
	for _, binary := range scenario.ExpectedBinaries {
		assert.True(t, tc.CheckBinaryExists(t, binary),
			"Binary %s should exist after installation", binary)
	}

	// Check config files exist
	for _, config := range scenario.ExpectedConfigs {
		assert.True(t, tc.CheckConfigExists(t, config),
			"Config file %s should exist after installation", config)
	}

	// Run validation commands
	for _, validationCmd := range scenario.ValidationCommands {
		output, err := tc.ExecuteCommand(t, validationCmd.Command, validationCmd.Args...)

		if validationCmd.ShouldFail {
			assert.Error(t, err, "Command %s should have failed", validationCmd.Command)
		} else {
			assert.NoError(t, err, "Command %s should have succeeded", validationCmd.Command)

			if validationCmd.ExpectedOut != "" {
				assert.Contains(t, output, validationCmd.ExpectedOut,
					"Command output should contain expected string")
			}
		}
	}
}

func (tc *TestContainer) Cleanup(t *testing.T) {
	if tc.DockerCmd != nil && tc.DockerCmd.Process != nil {
		tc.DockerCmd.Process.Kill()
	}

	// Remove container
	exec.Command("docker", "rm", "-f", tc.Name).Run()
}

func getProjectRoot() string {
	wd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	return wd
}
