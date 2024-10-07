package main

import (
	"github.com/StandardRunbook/plugin-interface/shared"
	plugin2 "github.com/StandardRunbook/plugin-shell-script/pkg/plugin"
	"github.com/hashicorp/go-plugin"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"shell-script-plugin": &shared.GRPCPlugin{Impl: &plugin2.ShellScriptPlugin{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
