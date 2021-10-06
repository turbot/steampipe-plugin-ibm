package ibm

import (
	"context"
	"strings"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// BuildServiceInstanceList :: return a list of matrix items, one per service instance
func BuildServiceInstanceList(ctx context.Context, connection *plugin.Connection) []map[string]interface{} {
	// cache service instance region matrix
	cacheKey := "ServiceInstanceList"

	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	// get all the service instances in the account
	serviceInstances, err := listAllServiceInstances(ctx, pluginQueryData, connection)
	if err != nil {
		panic(err)
	}

	matrix := make([]map[string]interface{}, len(serviceInstances))
	for i, instance := range serviceInstances {
		splitID := strings.Split(*instance, ":")
		matrix[i] = map[string]interface{}{
			"instance_id":  splitID[7],
			"instance_crn": *instance,
			"region":       splitID[5],
			"service_type": splitID[4],
		}
	}

	// set ServiceInstanceList cache
	pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)

	return matrix
}

func listAllServiceInstances(ctx context.Context, d *plugin.QueryData, connection *plugin.Connection) ([]*string, error) {
	// Create Session
	pluginQueryData.Connection = connection
	session, err := resourceControllerService(ctx, pluginQueryData)
	if err != nil {
		plugin.Logger(ctx).Error("listAllServiceInstances", "connection_error", err)
		return nil, err
	}

	serviceCacheKey := "listAllServiceInstances"
	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.([]*string), nil
	}

	var serviceInstanceCRNs []*string

	opts := &resourcecontrollerv2.ListResourceInstancesOptions{
		Type: core.StringPtr("service_instance"),
	}

	response, _, err := session.ListResourceInstances(opts)
	if err != nil {
		plugin.Logger(ctx).Error("listAllServiceInstances", "query_error", err)
		return nil, err
	}

	for _, i := range response.Resources {
		serviceInstanceCRNs = append(serviceInstanceCRNs, i.CRN)
	}

	// save service instances in cache
	pluginQueryData.ConnectionManager.Cache.Set(serviceCacheKey, serviceInstanceCRNs)

	return serviceInstanceCRNs, err
}
