package ibm

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v4/connection"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

// Regions is the current known list of valid regions
func Regions() []string {
	return []string{
		"au-syd",
		"br-sao",
		"ca-tor",
		"eu-de",
		"eu-gb",
		"jp-osa",
		"jp-tok",
		"us-east",
		"us-south",
	}
}

var pluginQueryData *plugin.QueryData

func init() {
	pluginQueryData = &plugin.QueryData{
		ConnectionManager: connection.NewManager(),
	}
}

// BuildRegionList :: return a list of matrix items, one per region specified in the connection config
func BuildRegionList(_ context.Context, connection *plugin.Connection) []map[string]interface{} {
	pluginQueryData.Connection = connection

	// cache matrix
	cacheKey := "RegionListMatrix"
	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	var allRegions []string

	// retrieve regions from connection config
	ibmConfig := GetConfig(connection)
	if ibmConfig.Regions != nil {
		regions := Regions()
		for _, pattern := range ibmConfig.Regions {
			for _, validRegion := range regions {
				if ok, _ := path.Match(pattern, validRegion); ok {
					allRegions = append(allRegions, validRegion)
				}
			}
		}
	}

	// Build regions matrix using config regions
	if len(allRegions) > 0 {
		uniqueRegions := unique(allRegions)

		if len(getInvalidRegions(uniqueRegions)) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(getInvalidRegions(uniqueRegions), ","))
		}

		// validate regions list
		matrix := make([]map[string]interface{}, len(uniqueRegions))
		for i, region := range uniqueRegions {
			matrix[i] = map[string]interface{}{"region": region}
		}

		// set cache
		pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)

		return matrix
	}

	// Search for region configured using env, or use default region (i.e. us-south)
	defaultIBMRegion := GetDefaultIBMRegion(pluginQueryData)
	matrix := []map[string]interface{}{
		{"region": defaultIBMRegion},
	}

	// set cache
	pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)

	return matrix
}

// Return invalid regions from a region list
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

// Transform used to get the instance_id column
func getServiceInstanceID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	instanceID := plugin.GetMatrixItem(ctx)["instance_id"].(string)
	return instanceID, nil
}

// Transform used to get the instance_crn column
func getServiceInstanceCRN(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	instanceID := plugin.GetMatrixItem(ctx)["instance_crn"].(string)
	return instanceID, nil
}

// GetDefaultIBMRegion returns the default region for IBM account
// if not set by Env variable
func GetDefaultIBMRegion(d *plugin.QueryData) string {
	// have we already created and cached the service?
	serviceCacheKey := "GetDefaultIBMRegion"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(string)
	}
	allIBMRegions := Regions()

	// get ibm config info
	ibmConfig := GetConfig(d.Connection)

	var regions []string
	var region string

	if ibmConfig.Regions != nil {
		regions = ibmConfig.Regions
		region = regions[0]
	} else {
		// Fetch regions from environment variables
		if os.Getenv("IBMCLOUD_REGION") != "" {
			region = os.Getenv("IBMCLOUD_REGION")
		}
		if os.Getenv("IC_REGION") != "" {
			region = os.Getenv("IC_REGION")
		}

		if region != "" {
			regions = []string{region}
		}

		// https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs#region
		if !helpers.StringSliceContains(allIBMRegions, region) {
			regions = []string{"us-south"}
		}
	}

	validPatterns := []string{}
	invalidPatterns := []string{}
	for _, namePattern := range regions {
		validRegions := []string{}
		for _, validRegion := range allIBMRegions {
			if ok, _ := path.Match(namePattern, validRegion); ok {
				validRegions = append(validRegions, validRegion)
			}
		}
		if len(validRegions) == 0 {
			invalidPatterns = append(invalidPatterns, namePattern)
		} else {
			validPatterns = append(validPatterns, namePattern)
		}
	}

	if len(validPatterns) == 0 {
		panic("\nconnection config have invalid \"regions\": " + strings.Join(invalidPatterns, ", ") + ". Edit your connection configuration file and then restart Steampipe")
	}

	if !helpers.StringSliceContains(allIBMRegions, region) {
		region = "us-south"
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, region)
	return region
}

// Returns a list of unique items
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
