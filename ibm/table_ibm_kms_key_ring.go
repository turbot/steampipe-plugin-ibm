package ibm

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmKmsKeyRing(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_kms_key_ring",
		Description:   "A key ring is a collection of leys in an IBM cloud location.",
		GetMatrixItemFunc: BuildServiceInstanceList,
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
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The key protect instance GUID.", Hydrate: plugin.HydrateFunc(getServiceInstanceID)},
			{Name: "creation_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the key ring was created."},
			{Name: "created_by", Type: proto.ColumnType_STRING, Description: "The creator of the key ring."},

			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getAccountId).WithCache(), Transform: transform.FromValue(), Description: "The account ID of this key ring."},
			{Name: "region", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getRegion), Description: "The region of this key ring."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: resourceInterfaceDescription("title")},
		},
	}
}

//// LIST FUNCTION

func listKmsKeyRings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listKmsKeyRings")

	instanceID := d.EqualsQualString("instance_id")
	serviceType := d.EqualsQualString("service_type")

	// Invalid service type
	if serviceType != "kms" {
		return nil, nil
	}

	// Return if specified instanceID not matched
	if d.EqualsQuals["instance_id"] != nil && d.EqualsQuals["instance_id"].GetStringValue() != instanceID {
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
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
