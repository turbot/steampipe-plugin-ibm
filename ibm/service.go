package ibm

import (
	"context"
	"fmt"

	"github.com/IBM/go-sdk-core/v5/core"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// kmsService return the service for IBM KMS service
func kmsService(ctx context.Context, d *plugin.QueryData) (*kp.Client, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "ibm_kms"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*kp.Client), nil
	}

	// Create region endpoint
	endpoint := fmt.Sprintf("https://%s.kms.cloud.ibm.com", region)
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	opts := kp.ClientConfig{
		BaseURL:  endpoint,
		APIKey:   apiKey,
		TokenURL: kp.DefaultTokenURL,
	}
	service, err := kp.New(opts, kp.DefaultTransport())
	if err != nil {
		return nil, err
	}
	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, service)
	return service, nil
}

// vpcService returns the service for IBM VPC Infrastructure service
func vpcService(ctx context.Context, d *plugin.QueryData, region string) (*vpcv1.VpcV1, error) {
	if region == "" {
		return nil, fmt.Errorf("region must be passed vpcService")
	}

	// Create region endpoint
	endpoint := fmt.Sprintf("https://%s.iaas.cloud.ibm.com/v1", region)

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := endpoint
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*vpcv1.VpcV1), nil
	}

	// Fetch API key from config
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

func iamService(ctx context.Context, d *plugin.QueryData) (*iamidentityv1.IamIdentityV1, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "ibm_iam"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*iamidentityv1.IamIdentityV1), nil
	}
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	serviceClientOptions := &iamidentityv1.IamIdentityV1Options{Authenticator: &core.IamAuthenticator{
		ApiKey: apiKey,
	}}
	// Instantiate the service with an API key based IAM authenticator
	service, err := iamidentityv1.NewIamIdentityV1UsingExternalConfig(serviceClientOptions)
	if err != nil {
		return nil, err
	}
	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, service)
	return service, nil
}

func iamAccessGroupService(ctx context.Context, d *plugin.QueryData) (*iamaccessgroupsv2.IamAccessGroupsV2, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "ibm_access_group"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*iamaccessgroupsv2.IamAccessGroupsV2), nil
	}
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	serviceClientOptions := &iamaccessgroupsv2.IamAccessGroupsV2Options{Authenticator: &core.IamAuthenticator{
		ApiKey: apiKey,
	}}
	service, err := iamaccessgroupsv2.NewIamAccessGroupsV2UsingExternalConfig(serviceClientOptions)
	if err != nil {
		return nil, err
	}
	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, service)
	return service, nil
}

func iamUserPolicy(ctx context.Context, d *plugin.QueryData) (*iampolicymanagementv1.IamPolicyManagementV1, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "ibm_user_policy"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*iampolicymanagementv1.IamPolicyManagementV1), nil
	}
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	serviceClientOptions := &iampolicymanagementv1.IamPolicyManagementV1Options{Authenticator: &core.IamAuthenticator{
		ApiKey: apiKey,
	}}
	service, err := iampolicymanagementv1.NewIamPolicyManagementV1UsingExternalConfig(serviceClientOptions)
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

func resourceControllerService(ctx context.Context, d *plugin.QueryData) (*resourcecontrollerv2.ResourceControllerV2, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "ibm_resource_controller"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*resourcecontrollerv2.ResourceControllerV2), nil
	}
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	opts := &resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	}
	service, err := resourcecontrollerv2.NewResourceControllerV2(opts)
	if err != nil {
		return nil, err
	}
	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, service)
	return service, nil
}

func resourceManagerService(ctx context.Context, d *plugin.QueryData) (*resourcemanagerv2.ResourceManagerV2, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "ibm_resource_manager"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*resourcemanagerv2.ResourceManagerV2), nil
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
