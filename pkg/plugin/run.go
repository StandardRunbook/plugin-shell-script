package plugin

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	pluginInterface "github.com/StandardRunbook/plugin-interface"
	"github.com/StandardRunbook/plugin-shell-script/pkg/config"
)

//go:embed run.sh
var runScript []byte

type ShellScriptConfig struct {
	name           string
	version        string
	arguments      []string
	output         string
	expectedOutput string
}

func (t *ShellScriptConfig) Name() string {
	return t.name
}

func (t *ShellScriptConfig) Version() string {
	return t.version
}

func (t *ShellScriptConfig) Run() error {
	// Step 1: Create a temporary file
	tmpFile, err := os.CreateTemp("", "script-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Ensure the file is removed after execution

	// Step 2: Write the embedded script to the temporary file
	_, err = tmpFile.Write(runScript)
	if err != nil {
		return fmt.Errorf("failed to write '%s' script to temporary file: %w", t.Name(), err)
	}

	// Step 3: Close the file to flush writes and prepare it for execution
	err = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Step 4: Set the appropriate permissions to make the script executable
	err = os.Chmod(tmpFile.Name(), 0755)
	if err != nil {
		return fmt.Errorf("failed to set executable permissions on file: %w", err)
	}

	if len(t.arguments) == 0 {
		return fmt.Errorf("no file path provided")
	}

	// Step 5: Execute the script
	cmd := exec.Command("/bin/bash", tmpFile.Name(), t.arguments[0])
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing script: %w", err)
	}
	t.output = string(output)
	return nil
}

func (t *ShellScriptConfig) ParseOutput() string {
	if strings.Contains(t.output, t.expectedOutput) {
		return "success"
	}
	return "failure"
}

func NewShellScriptPlugin(cfg *config.ShellScriptConfig) pluginInterface.IPlugin {
	return &ShellScriptConfig{
		name:           cfg.Name,
		version:        cfg.Version,
		arguments:      cfg.ScriptArguments,
		expectedOutput: cfg.ExpectedOutput,
	}
}
