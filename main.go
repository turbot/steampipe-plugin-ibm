package main

import (
	"github.com/turbot/steampipe-plugin-ibm/ibm"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: ibm.Plugin})
}
