package ibm

import (
	"context"
	"fmt"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/networking-go-sdk/globalloadbalancerv1"
	"github.com/IBM/networking-go-sdk/zonessettingsv1"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
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

// cisZoneService returns the service for IBM CIS Zone service
func cisZoneService(ctx context.Context, d *plugin.QueryData) (*zonesv1.ZonesV1, error) {
	serviceInstanceID := "crn:v1:bluemix:public:internet-svcs:global:a/76aa4877fab6436db86f121f62faf221:3e5dc1e0-3aea-4699-986e-e3b8f117c51d::"
	endpoint := "https://api.cis.cloud.ibm.com"

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "cisZone"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*zonesv1.ZonesV1), nil
	}

	// Fetch API key from config
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	// Instantiate the service with an API key based IAM authenticator
	service, err := zonesv1.NewZonesV1(&zonesv1.ZonesV1Options{
		Crn: &serviceInstanceID,
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

// cisZoneSettingService returns the service for IBM CIS Zone Setting service
func cisZoneSettingService(ctx context.Context, d *plugin.QueryData, zoneId string) (*zonessettingsv1.ZonesSettingsV1, error) {
	serviceInstanceID := "crn:v1:bluemix:public:internet-svcs:global:a/76aa4877fab6436db86f121f62faf221:3e5dc1e0-3aea-4699-986e-e3b8f117c51d::"
	endpoint := "https://api.cis.cloud.ibm.com"

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "cisZoneSetting"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*zonessettingsv1.ZonesSettingsV1), nil
	}

	// Fetch API key from config
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	// Instantiate the service with an API key based IAM authenticator
	service, err := zonessettingsv1.NewZonesSettingsV1(&zonessettingsv1.ZonesSettingsV1Options{
		Crn:            &serviceInstanceID,
		ZoneIdentifier: &zoneId,
		URL:            endpoint,
		Authenticator: &core.IamAuthenticator{
			ApiKey: apiKey,
		},
	})
	if err != nil {
		plugin.Logger(ctx).Error("getTlsMinimumVersion", "zoneId", zoneId)
		return nil, err
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, service)

	return service, nil
}

// cisGlobalLoadBalancerService returns the service for IBM CIS Global Load Balancer service
func cisGlobalLoadBalancerService(ctx context.Context, d *plugin.QueryData, zoneId string) (*globalloadbalancerv1.GlobalLoadBalancerV1, error) {
	serviceInstanceID := "crn:v1:bluemix:public:internet-svcs:global:a/76aa4877fab6436db86f121f62faf221:3e5dc1e0-3aea-4699-986e-e3b8f117c51d::"
	endpoint := "https://api.cis.cloud.ibm.com"

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "cisGlobalLoadBalancer"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*globalloadbalancerv1.GlobalLoadBalancerV1), nil
	}

	// Fetch API key from config
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	// Instantiate the service with an API key based IAM authenticator
	service, err := globalloadbalancerv1.NewGlobalLoadBalancerV1(&globalloadbalancerv1.GlobalLoadBalancerV1Options{
		Crn:            &serviceInstanceID,
		ZoneIdentifier: &zoneId,
		URL:            endpoint,
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

// cisGlobalLoadBalancerService returns the service for IBM CIS DNS service
func cisDnsRecordService(ctx context.Context, d *plugin.QueryData, zoneId string) (*dnsrecordsv1.DnsRecordsV1, error) {
	serviceInstanceID := "crn:v1:bluemix:public:internet-svcs:global:a/76aa4877fab6436db86f121f62faf221:3e5dc1e0-3aea-4699-986e-e3b8f117c51d::"
	endpoint := "https://api.cis.cloud.ibm.com"

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "cisDnsRecord"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*dnsrecordsv1.DnsRecordsV1), nil
	}

	// Fetch API key from config
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	// Instantiate the service with an API key based IAM authenticator
	service, err := dnsrecordsv1.NewDnsRecordsV1(&dnsrecordsv1.DnsRecordsV1Options{
		Crn:            &serviceInstanceID,
		ZoneIdentifier: &zoneId,
		URL:            endpoint,
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

func iamPolicyManagementService(ctx context.Context, d *plugin.QueryData) (*iampolicymanagementv1.IamPolicyManagementV1, error) {
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

func cosService(ctx context.Context, d *plugin.QueryData, region string) (*s3.S3, error) {

	serviceInstanceID := plugin.GetMatrixItem(ctx)["instance_crn"].(string)
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "ibm_cos" + region
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*s3.S3), nil
	}
	apiKey, err := configApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	authEndpoint := "https://iam.cloud.ibm.com/identity/token"
	serviceEndpoint := fmt.Sprintf("s3.%s.cloud-object-storage.appdomain.cloud", region)

	conf := aws.NewConfig().
		WithEndpoint(serviceEndpoint).
		WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(),
			authEndpoint, apiKey, serviceInstanceID)).
		WithS3ForcePathStyle(true)

	sess := session.Must(session.NewSession())

	service := s3.New(sess, conf)

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, service)
	return service, nil
}
