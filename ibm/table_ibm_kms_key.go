package ibm

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	kp "github.com/IBM/keyprotect-go-client"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmKmsKey(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_kms_key",
		Description:   "Users in the IBM Cloud account.",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listKmsKeys,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "instance_id",
					Require: plugin.Required,
				},
				{
					Name:    "key_ring_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getKmsKey,
			KeyColumns: plugin.AllColumns([]string{"instance_id", "id"}),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "A human-readable name assigned to your key for convenience."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "An unique identifier of the key."},
			{Name: "crn", Type: proto.ColumnType_STRING, Description: "The Cloud Resource Name (CRN) that uniquely identifies your cloud resources.", Transform: transform.FromField("CRN")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Specifies the MIME type that represents the key resource."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The key state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1, Suspended = 2, Deactivated = 3, and Destroyed = 5 values."},
			{Name: "imported", Type: proto.ColumnType_BOOL, Description: "Indicates whether the key was originally imported or generated in Key Protect."},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The key protect instance GUID.", Transform: transform.FromQual("instance_id")},
			{Name: "algorithm_type", Type: proto.ColumnType_STRING, Description: "Specifies the key algorithm."},
			{Name: "creation_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date the key material was created."},
			{Name: "created_by", Type: proto.ColumnType_STRING, Description: "The unique identifier for the resource that created the key."},
			{Name: "deleted", Type: proto.ColumnType_BOOL, Description: "Indicates whether the key has been deleted, or not."},
			{Name: "deleted_by", Type: proto.ColumnType_STRING, Description: "The unique identifier for the resource that deleted the key."},
			{Name: "deletion_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date the key material was destroyed."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A text field used to provide a more detailed description of the key."},
			{Name: "encrypted_nonce", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "encryption_algorithm", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "expiration", Type: proto.ColumnType_TIMESTAMP, Description: "The date the key material will expire."},
			{Name: "extractable", Type: proto.ColumnType_BOOL, Description: "Indicates whether the key material can leave the service, or not."},
			{Name: "key_ring_id", Type: proto.ColumnType_STRING, Description: "An ID that identifies the key ring.", Transform: transform.FromField("KeyRingID")},
			{Name: "last_rotate_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date when the key was last rotated."},
			{Name: "last_update_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date when the key metadata was last modified."},
			{Name: "payload", Type: proto.ColumnType_STRING, Description: "Specifies the key payload."},
			{Name: "aliases", Type: proto.ColumnType_JSON, Description: "A list of key aliases."},
			{Name: "dual_auth_delete", Type: proto.ColumnType_JSON, Description: "Metadata that indicates the status of a dual authorization policy on the key."},
			{Name: "key_version", Type: proto.ColumnType_JSON, Description: "Properties associated with a specific key version."},

			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN").Transform(crnToAccountID), Description: "The account ID of this key."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this key."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRN").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: resourceInterfaceDescription("tags")},
		},
	}
}

//// LIST FUNCTION

func listKmsKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)
	instanceID := d.KeyColumnQuals["instance_id"].GetStringValue()

	// Get service instance details
	instanceData, err := getServiceInstance(ctx, d, instanceID)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_kms_key.getServiceInstance", "query_error", err)
	}
	splitID := strings.Split(*instanceData.ID, ":")
	var instanceRegion string
	if len(splitID) > 5 {
		instanceRegion = splitID[5]
	}

	// Compare service instance region with config region
	if region != instanceRegion {
		return nil, nil
	}

	endpointType := d.KeyColumnQuals["endpoint_type"].GetStringValue()

	// Create service connection
	conn, err := kmsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_kms_key.listKmsKeys", "connection_error", err)
		return nil, err
	}
	conn.Config.InstanceID = instanceID

	URL, err := KmsEndpointURL(conn, endpointType, instanceData.Extensions)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_kms_key.KmsEndpointURL", "connection_error", err)
		return nil, err
	}
	conn.URL = URL

	// Additional filters
	if d.KeyColumnQuals["key_ring_id"] != nil {
		conn.Config.KeyRing = d.KeyColumnQuals["key_ring_id"].GetStringValue()
	}

	// Retrieve the list of keys for your account.
	maxResult := int64(100)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < maxResult {
			maxResult = *limit
		}
	}

	data, err := conn.GetKeys(ctx, int(maxResult), 0)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_kms_key.listKmsKeys", "query_error", err)
		return nil, err
	}
	for _, i := range data.Keys {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getKmsKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)
	instanceID := d.KeyColumnQuals["instance_id"].GetStringValue()
	id := d.KeyColumnQuals["id"].GetStringValue()

	// No inputs
	if instanceID == "" && id == "" {
		return nil, nil
	}

	// Get service instance details
	instanceData, err := getServiceInstance(ctx, d, instanceID)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_kms_key.getServiceInstance", "query_error", err)
	}
	splitID := strings.Split(*instanceData.ID, ":")
	var instanceRegion string
	if len(splitID) > 5 {
		instanceRegion = splitID[5]
	}

	// Compare service instance region with config region
	if region != instanceRegion {
		return nil, nil
	}

	// Create service connection
	conn, err := kmsService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_kms_key.listKmsKeys", "connection_error", err)
		return nil, err
	}
	conn.Config.InstanceID = instanceID

	data, err := conn.GetKeyMetadata(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_kms_key.listKmsKeys", "query_error", err)
		return nil, err
	}

	return data, nil
}

func KmsEndpointURL(kpAPI *kp.Client, endpointType string, extensions map[string]interface{}) (*url.URL, error) {
	exturl := extensions["endpoints"].(map[string]interface{})["public"]
	if endpointType == "private" || strings.Contains(kpAPI.Config.BaseURL, "private") {
		exturl = extensions["endpoints"].(map[string]interface{})["private"]
	}
	endpointURL := fmt.Sprintf("%s/api/v2/keys", exturl.(string))

	u, err := url.Parse(endpointURL)
	if err != nil {
		return nil, err
	}
	return u, nil
}
