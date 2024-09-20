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
	require.Equal(t, shellScriptPlugin.Name(), "test-shell-script")
	require.Equal(t, shellScriptPlugin.Version(), "v1.0.0")
	err := shellScriptPlugin.Run()
	require.NoError(t, err)
	require.Equal(t, shellScriptPlugin.ParseOutput(), "success")
}
