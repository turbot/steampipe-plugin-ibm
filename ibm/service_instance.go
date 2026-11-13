package ibm

import (
	"context"
	"strings"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

// BuildServiceInstanceList :: return a list of matrix items, one per service instance
func BuildServiceInstanceList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	// cache service instance region matrix
	cacheKey := "ServiceInstanceList"

	if cachedData, ok := d.ConnectionCache.Get(ctx, cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	// get all the service instances in the account
	serviceInstances, err := listAllServiceInstances(ctx, d, d.Connection)
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
	_ = d.ConnectionCache.Set(ctx, cacheKey, matrix)

	return matrix
}

func listAllServiceInstances(ctx context.Context, d *plugin.QueryData, connection *plugin.Connection) ([]*string, error) {
	// Create Session
	session, err := resourceControllerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listAllServiceInstances", "connection_error", err)
		return nil, err
	}

	serviceCacheKey := "listAllServiceInstances"
	if cachedData, ok := d.ConnectionCache.Get(ctx, serviceCacheKey); ok {
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
	_ = d.ConnectionCache.Set(ctx, serviceCacheKey, serviceInstanceCRNs)

	return serviceInstanceCRNs, err
}
