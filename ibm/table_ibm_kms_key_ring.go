package ibm

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmKmsKeyRing(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_kms_key_ring",
		Description:   "A key ring is a collection of leys in an IBM cloud location.",
		GetMatrixItem: BuildServiceInstanceList,
		List: &plugin.ListConfig{
			Hydrate: listKmsKeyRings,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "instance_id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "An unique identifier of the key ring."},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The key protect instance GUID.", Transform: transform.From(getServiceInstanceID)},
			{Name: "creation_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the key ring was created."},
			{Name: "created_by", Type: proto.ColumnType_STRING, Description: "The creator of the key ring."},

			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getAccountId).WithCache(), Transform: transform.FromValue(), Description: "The account ID of this key ring."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this key ring."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: resourceInterfaceDescription("title")},
		},
	}
}

//// LIST FUNCTION

func listKmsKeyRings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKmsKeyRings")

	instanceID := plugin.GetMatrixItem(ctx)["instance_id"].(string)
	serviceType := plugin.GetMatrixItem(ctx)["service_type"].(string)

	// Invalid service type
	if serviceType != "kms" {
		return nil, nil
	}

	// Return if specified instanceID not matched
	if d.KeyColumnQuals["instance_id"] != nil && d.KeyColumnQuals["instance_id"].GetStringValue() != instanceID {
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
