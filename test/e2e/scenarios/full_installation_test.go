package scenarios

import (
	"testing"

	"github.com/petersenjoern/devenv/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestFullDeveloperWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	scenarios := []struct {
		name        string
		environment string
		workflow    []WorkflowStep
	}{
		{
			name:        "Complete Developer Setup - Ubuntu",
			environment: "ubuntu",
			workflow: []WorkflowStep{
				{
					Name:        "Initial Status Check",
					Command:     []string{"status"},
					ExpectError: false,
				},
				{
					Name:        "Install Core Tools",
					Command:     []string{"install", "--non-interactive", "--tools", "zsh,vim,git"},
					ExpectError: false,
				},
				{
					Name:        "Verify Installation",
					Command:     []string{"status"},
					ExpectError: false,
					Validations: []string{"zsh: installed", "vim: installed", "git: installed"},
				},
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			container := integration.SetupTestContainer(t, scenario.environment)
			defer container.Cleanup(t)

			for _, step := range scenario.workflow {
				t.Run(step.Name, func(t *testing.T) {
					output, err := container.ExecuteDevEnv(t, step.Command...)

					if step.ExpectError {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
					}

					for _, validation := range step.Validations {
						assert.Contains(t, output, validation)
					}
				})
			}
		})
	}
}

type WorkflowStep struct {
	Name        string
	Command     []string
	ExpectError bool
	Validations []string
}
