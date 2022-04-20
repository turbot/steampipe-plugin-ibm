package ibm

import (
	"context"
	"strings"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableCosBucket(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_cos_bucket",
		Description:   "An IBM Cloud storage bucket.",
		GetMatrixItem: BuildServiceInstanceList,
		List: &plugin.ListConfig{
			Hydrate: listBucket,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the bucket."},
			{Name: "creation_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date when the bucket was created."},
			{Name: "sse_kp_customer_root_key_crn", Type: proto.ColumnType_STRING, Description: "The root key used by Key Protect to encrypt this bucket. This value must be the full CRN of the root key.", Hydrate: headBucket, Transform: transform.FromField("IBMSSEKPCrkId")},
			{Name: "sse_kp_enabled", Type: proto.ColumnType_BOOL, Description: "Specifies whether the Bucket has Key Protect enabled.", Hydrate: headBucket, Transform: transform.FromField("IBMSSEKPEnabled"), Default: false},
			{Name: "versioning_enabled", Type: proto.ColumnType_BOOL, Description: "The versioning state of a bucket.", Hydrate: getBucketVersioning, Transform: transform.FromField("Status").Transform(handleNilString).Transform(transform.ToBool)},
			{Name: "versioning_mfa_delete", Type: proto.ColumnType_BOOL, Description: "The MFA Delete status of the versioning state.", Hydrate: getBucketVersioning, Transform: transform.FromField("MFADelete").Transform(handleNilString).Transform(transform.ToBool)},
			{Name: "acl", Type: proto.ColumnType_JSON, Description: "The access control list (ACL) of a bucket.", Hydrate: getBucketACL, Transform: transform.FromValue()},
			{Name: "lifecycle_rules", Type: proto.ColumnType_JSON, Description: "The lifecycle configuration information of the bucket.", Hydrate: getBucketLifecycle, Transform: transform.FromField("Rules")},
			{Name: "public_access_block_configuration", Type: proto.ColumnType_JSON, Description: "The public access block configuration information of the bucket.", Hydrate: getBucketPublicAccessBlockConfiguration, Transform: transform.FromValue()},
			{Name: "retention", Type: proto.ColumnType_JSON, Description: "The retention configuration information of the bucket.", Hydrate: getBucketRetention, Transform: transform.FromValue()},
			{Name: "website", Type: proto.ColumnType_JSON, Description: "The lifecycle configuration information of the bucket.", Hydrate: getBucketWebsite, Transform: transform.FromValue()},

			// Standard columns
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.FromField("LocationConstraint"), Description: "The region of the bucket."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
		},
	}
}

//// LIST FUNCTION

func listBucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := GetDefaultIBMRegion(d)
	plugin.Logger(ctx).Trace("listBucket")

	serviceType := plugin.GetMatrixItem(ctx)["service_type"].(string)

	// Invalid service type
	if serviceType != "cloud-object-storage" {
		return nil, nil
	}

	// Create service connection
	conn, err := cosService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cos_bucket.listBucket", "connection_error", err)
		return nil, err
	}

	opt := &s3.ListBucketsExtendedInput{}

	data, err := conn.ListBucketsExtended(opt)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cos_bucket.listBucket", "query_error", err)
		return nil, err
	}
	for _, i := range data.Buckets {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func headBucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("headBucket")
	serviceType := plugin.GetMatrixItem(ctx)["service_type"].(string)

	// Invalid service type
	if serviceType != "cloud-object-storage" {
		return nil, nil
	}
	bucket := h.Item.(*s3.BucketExtended)

	location := strings.TrimSuffix(*bucket.LocationConstraint, "-smart")
	// Create Session
	conn, err := cosService(ctx, d, location)
	if err != nil {
		return nil, err
	}

	params := &s3.HeadBucketInput{
		Bucket: bucket.Name,
	}

	lifecycle, err := conn.HeadBucket(params)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cos_bucket.headBucket", "query_error", err)
		if strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		}
		return nil, err
	}
	return lifecycle, nil
}

func getBucketLifecycle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketLifecycle")
	serviceType := plugin.GetMatrixItem(ctx)["service_type"].(string)

	// Invalid service type
	if serviceType != "cloud-object-storage" {
		return nil, nil
	}
	bucket := h.Item.(*s3.BucketExtended)

	location := strings.TrimSuffix(*bucket.LocationConstraint, "-smart")
	// Create Session
	conn, err := cosService(ctx, d, location)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketLifecycleConfigurationInput{
		Bucket: bucket.Name,
	}

	lifecycle, err := conn.GetBucketLifecycleConfiguration(params)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cos_bucket.getBucketLifecycle", "query_error", err)
		if strings.Contains(err.Error(), "lifecycle configuration does not exist") {
			return nil, nil
		}
		return nil, err
	}
	return lifecycle, nil
}

func getBucketRetention(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketRetention")
	serviceType := plugin.GetMatrixItem(ctx)["service_type"].(string)

	// Invalid service type
	if serviceType != "cloud-object-storage" {
		return nil, nil
	}
	bucket := h.Item.(*s3.BucketExtended)

	location := strings.TrimSuffix(*bucket.LocationConstraint, "-smart")
	// Create Session
	conn, err := cosService(ctx, d, location)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketProtectionConfigurationInput{
		Bucket: bucket.Name,
	}

	retention, err := conn.GetBucketProtectionConfiguration(params)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cos_bucket.getBucketRetention", "query_error", err)
		if strings.Contains(err.Error(), "lifecycle configuration does not exist") {
			return nil, nil
		}
		return nil, err
	}
	return retention, nil
}

func getBucketVersioning(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketVersioning")

	serviceType := plugin.GetMatrixItem(ctx)["service_type"].(string)

	// Invalid service type
	if serviceType != "cloud-object-storage" {
		return nil, nil
	}
	bucket := h.Item.(*s3.BucketExtended)

	location := strings.TrimSuffix(*bucket.LocationConstraint, "-smart")

	// Create Session
	conn, err := cosService(ctx, d, location)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketVersioningInput{
		Bucket: bucket.Name,
	}

	versioning, err := conn.GetBucketVersioning(params)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cos_bucket.getBucketVersioning", "query_error", err)
		if strings.Contains(err.Error(), "bucket does not exist") {
			return nil, nil
		}
		return nil, err
	}

	return versioning, nil
}

func getBucketWebsite(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketWebsite")

	serviceType := plugin.GetMatrixItem(ctx)["service_type"].(string)

	// Invalid service type
	if serviceType != "cloud-object-storage" {
		return nil, nil
	}
	bucket := h.Item.(*s3.BucketExtended)

	location := strings.TrimSuffix(*bucket.LocationConstraint, "-smart")

	// Create Session
	conn, err := cosService(ctx, d, location)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketWebsiteInput{
		Bucket: bucket.Name,
	}

	website, err := conn.GetBucketWebsite(params)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cos_bucket.getBucketWebsite", "query_error", err)
		if strings.Contains(err.Error(), "bucket does not have a website configuration") {
			return nil, nil
		}
		return nil, err
	}

	return website, nil
}

func getBucketACL(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketACL")
	serviceType := plugin.GetMatrixItem(ctx)["service_type"].(string)

	// Invalid service type
	if serviceType != "cloud-object-storage" {
		return nil, nil
	}
	bucket := h.Item.(*s3.BucketExtended)

	location := strings.TrimSuffix(*bucket.LocationConstraint, "-smart")

	// Create Session
	conn, err := cosService(ctx, d, location)
	if err != nil {
		return nil, err
	}

	params := &s3.GetBucketAclInput{
		Bucket: bucket.Name,
	}

	acl, err := conn.GetBucketAcl(params)
	if err != nil {
		return nil, err
	}

	return acl, nil
}

func getBucketPublicAccessBlockConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBucketPublicAccessBlockConfiguration")
	serviceType := plugin.GetMatrixItem(ctx)["service_type"].(string)

	// Invalid service type
	if serviceType != "cloud-object-storage" {
		return nil, nil
	}
	bucket := h.Item.(*s3.BucketExtended)

	location := strings.TrimSuffix(*bucket.LocationConstraint, "-smart")

	// Create Session
	conn, err := cosService(ctx, d, location)
	if err != nil {
		return nil, err
	}

	params := &s3.GetPublicAccessBlockInput{
		Bucket: bucket.Name,
	}

	result, err := conn.GetPublicAccessBlock(params)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchPublicAccessBlockConfiguration") {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}
