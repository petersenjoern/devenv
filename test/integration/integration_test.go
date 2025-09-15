package integration

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	containers map[string]*TestContainer
}

func TestIntegrationSuite(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		t.Skip("Integration tests skipped. Set INTEGRATION_TESTS=1 to run")
	}

	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.containers = make(map[string]*TestContainer)

	// Setup containers for each environment
	environments := []string{"ubuntu"}

	for _, env := range environments {
		container := SetupTestContainer(suite.T(), env)
		suite.containers[env] = container
	}

	// Wait for containers to be ready
	time.Sleep(60 * time.Second)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	for _, container := range suite.containers {
		container.Cleanup(suite.T())
	}
}

func (suite *IntegrationTestSuite) TestFullInstallationUbuntu() {
	container := suite.containers["ubuntu"]

	scenario := TestScenario{
		Name:             "Full Installation on Ubuntu",
		Environment:      "ubuntu",
		ToolsToInstall:   []string{"zsh", "fzf", "neovim", "mise"},
		ExpectedBinaries: []string{"zsh", "fzf", "nvim", "mise"},
		ExpectedConfigs: []string{
			"/home/testuser/.zshrc",
			"/home/testuser/.config/nvim/init.vim",
			"/home/testuser/.config/mise/config.toml",
		},
		ValidationCommands: []ValidationCmd{
			{
				Command:     "zsh",
				Args:        []string{"--version"},
				ExpectedOut: "zsh",
			},
			{
				Command:     "nvim",
				Args:        []string{"--version"},
				ExpectedOut: "NVIM",
			},
		},
	}

	// Execute installation (simulate interactive selection)
	output, err := container.ExecuteDevEnv(suite.T(), "install", "--non-interactive", "--tools", "zsh,fzf,neovim,mise")
	suite.NoError(err, "DevEnv installation should succeed")
	suite.Contains(output, "Installation completed")

	// Validate installation
	container.ValidateInstallation(suite.T(), scenario)
}

func (suite *IntegrationTestSuite) TestPartialInstallationWSL() {
	container := suite.containers["wsl-ubuntu"]

	scenario := TestScenario{
		Name:             "Partial Installation on WSL",
		Environment:      "wsl-ubuntu",
		ToolsToInstall:   []string{"zsh", "fzf"},
		ExpectedBinaries: []string{"zsh", "fzf"},
		ExpectedConfigs: []string{
			"/home/testuser/.zshrc",
		},
	}

	output, err := container.ExecuteDevEnv(suite.T(), "install", "--non-interactive", "--tools", "zsh,fzf")
	suite.NoError(err)
	suite.Contains(output, "Installation completed")

	container.ValidateInstallation(suite.T(), scenario)
}

func (suite *IntegrationTestSuite) TestStatusCommand() {
	container := suite.containers["ubuntu"]

	// First install some tools
	container.ExecuteDevEnv(suite.T(), "install", "--non-interactive", "--tools", "zsh")

	// Then check status
	output, err := container.ExecuteDevEnv(suite.T(), "status")
	suite.NoError(err)
	suite.Contains(output, "zsh")
	suite.Contains(output, "installed") // or whatever status format you use
}
