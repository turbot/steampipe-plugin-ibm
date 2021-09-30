package ibm

import (
	"context"
	"fmt"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func vpcService(ctx context.Context, d *plugin.QueryData) (*vpcv1.VpcV1, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)
	endpoint := fmt.Sprintf("https://%s.iaas.cloud.ibm.com/v1", region)
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := endpoint
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*vpcv1.VpcV1), nil
	}
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	// Instantiate the service with an API key based IAM authenticator
	service, err := vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		URL: endpoint,
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	})
	if err != nil {
		return nil, err
	}
	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, service)
	return service, nil
}

func tagService(ctx context.Context, d *plugin.QueryData) (*globaltaggingv1.GlobalTaggingV1, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "ibm_global_tagging"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*globaltaggingv1.GlobalTaggingV1), nil
	}
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	opts := &globaltaggingv1.GlobalTaggingV1Options{
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	}
	service, err := globaltaggingv1.NewGlobalTaggingV1(opts)
	if err != nil {
		return nil, err
	}
	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, service)
	return service, nil
}

/*
func resourceManagerService(ctx context.Context, d *plugin.QueryData) (*resourcemanagerv2.ResourceManagementAPIv2, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "ibm_resource_controller"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*resourcemanagerv2.ResourceManagementAPIv2), nil
	}
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	opts := &resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	}
	service, err := resourcemanagerv2.NewResourceManagerV2(opts)
	if err != nil {
		return nil, err
	}
	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, service)
	return service, nil
}
*/
