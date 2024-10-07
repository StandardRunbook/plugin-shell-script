package plugin

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	pluginInterface "github.com/StandardRunbook/plugin-interface/shared"
	"github.com/StandardRunbook/plugin-shell-script/pkg/config"
)

//go:embed run.sh
var runScript []byte

type ShellScriptPlugin struct {
	name           string
	version        string
	arguments      []string
	output         string
	expectedOutput string
}

func (t *ShellScriptPlugin) Init(m map[string]string) error {
	t.name = m["name"]
	t.version = m["version"]
	t.arguments = strings.Split(m["arguments"], ",")
	t.expectedOutput = m["expected_output"]
	return nil
}

func (t *ShellScriptPlugin) Name() (string, error) {
	if strings.EqualFold(t.name, "") {
		return "invalid plugin", fmt.Errorf("plugin name is empty")
	}
	return t.name, nil
}

func (t *ShellScriptPlugin) Version() (string, error) {
	if strings.EqualFold(t.version, "") {
		return "invalid version", fmt.Errorf("plugin version is empty")
	}
	return t.version, nil
}

func (t *ShellScriptPlugin) Run() error {
	// Step 1: Create a temporary file
	tmpFile, err := os.CreateTemp("", "script-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Ensure the file is removed after execution

	// Step 2: Write the embedded script to the temporary file
	name, err := t.Name()
	if err != nil {
		return err
	}
	_, err = tmpFile.Write(runScript)
	if err != nil {
		return fmt.Errorf("failed to write '%s' script to temporary file: %w", name, err)
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

func (t *ShellScriptPlugin) ParseOutput() (string, error) {
	if strings.EqualFold(t.output, t.expectedOutput) {
		return "success", nil
	}
	return "failure", fmt.Errorf("output '%s' is not expected: %s", t.output, t.expectedOutput)
}

func NewShellScriptPlugin(cfg *config.ShellScriptConfig) pluginInterface.IPlugin {
	return &ShellScriptPlugin{
		name:           cfg.Name,
		version:        cfg.Version,
		arguments:      cfg.ScriptArguments,
		expectedOutput: cfg.ExpectedOutput,
	}
}
