package ibm

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmKmsKeyRing(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_kms_key_ring",
		Description:   "Users in the IBM Cloud account.",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listKmsKeyRings,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "instance_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "An unique identifier of the key ring."},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The key protect instance GUID.", Transform: transform.FromQual("instance_id")},
			{Name: "creation_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the key ring was created."},
			{Name: "created_by", Type: proto.ColumnType_STRING, Description: "The creator of the key ring."},
		},
	}
}

//// LIST FUNCTION

func listKmsKeyRings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)
	instanceID := d.KeyColumnQuals["instance_id"].GetStringValue()

	// Get service instance details
	instanceData, err := getServiceInstance(ctx, d, instanceID)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_kms_key_ring.getServiceInstance", "query_error", err)
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
		plugin.Logger(ctx).Error("ibm_kms_key_ring.listKmsKeyRings", "connection_error", err)
		return nil, err
	}
	conn.Config.InstanceID = instanceID

	data, err := conn.GetKeyRings(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_kms_key_ring.listKmsKeyRings", "query_error", err)
		return nil, err
	}
	for _, i := range data.KeyRings {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
