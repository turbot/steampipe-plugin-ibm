package ibm

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type ibmConfig struct {
	APIKey  *string  `cty:"api_key"`
	Regions []string `cty:"regions,optional"`
}

func ConfigInstance() interface{} {
	return &ibmConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) ibmConfig {
	if connection == nil || connection.Config == nil {
		return ibmConfig{}
	}
	config, _ := connection.Config.(ibmConfig)
	return config
}
