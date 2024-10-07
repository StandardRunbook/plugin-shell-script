package plugin_test

import (
	"testing"

	"github.com/StandardRunbook/plugin-shell-script/pkg/config"
	"github.com/StandardRunbook/plugin-shell-script/pkg/plugin"
	"github.com/stretchr/testify/require"
)

func TestTemplate_Run(t *testing.T) {
	t.Parallel()
	cfg := &config.ShellScriptConfig{
		Name:            "test-shell-script",
		Version:         "v1.0.0",
		ExpectedOutput:  "Hello World!",
		ScriptArguments: []string{"./test.sh"},
	}
	shellScriptPlugin := plugin.NewShellScriptPlugin(cfg)
	name, err := shellScriptPlugin.Name()
	require.NoError(t, err)
	require.Equal(t, name, "test-shell-script")
	version, err := shellScriptPlugin.Version()
	require.NoError(t, err)
	require.Equal(t, version, "v1.0.0")
	err = shellScriptPlugin.Run()
	require.NoError(t, err)
	output, err := shellScriptPlugin.ParseOutput()
	require.NoError(t, err)
	require.Equal(t, output, "success")
}
