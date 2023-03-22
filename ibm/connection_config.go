package ibm

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type ibmConfig struct {
	APIKey  *string  `cty:"api_key"`
	Regions []string `cty:"regions"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"api_key": {
		Type: schema.TypeString,
	},
	"regions": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
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
