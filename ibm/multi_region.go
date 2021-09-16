package ibm

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// Regions is the current known list of valid regions
func Regions() []string {
	return []string{
		"au-syd",
		"eu-de",
		"eu-gb",
		"jp-osa",
		"jp-tok",
		"us-east",
		"us-south",
	}
}

// DefaultRegions is the list of regions used in queries by default
func DefaultRegions() []string {
	return []string{
		"us-south",
	}
}

// BuildRegionList :: return a list of matrix items, one per region specified in the connection config
func BuildRegionList(_ context.Context, connection *plugin.Connection) []map[string]interface{} {
	// retrieve regions from connection config
	ibmConfig := GetConfig(connection)
	if &ibmConfig != nil && ibmConfig.Regions != nil {
		regions := GetConfig(connection).Regions

		if len(getInvalidRegions(regions)) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(getInvalidRegions(regions), ","))
		}

		// validate regions list
		matrix := make([]map[string]interface{}, len(regions))
		for i, region := range regions {
			matrix[i] = map[string]interface{}{"region": region}
		}
		return matrix
	}

	return []map[string]interface{}{
		//TODO
		{"region": "us-south"},
	}
}

func getInvalidRegions(regions []string) []string {
	invalidRegions := []string{}
	for _, region := range regions {
		if !helpers.StringSliceContains(Regions(), region) {
			invalidRegions = append(invalidRegions, region)
		}
	}
	return invalidRegions
}

// Transform used to get the region column
func getRegion(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)
	return region, nil
}
