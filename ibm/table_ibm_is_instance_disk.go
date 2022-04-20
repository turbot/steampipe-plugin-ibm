package ibm

import (
	"context"
	"strings"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmIsInstanceDisk(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_is_instance_disk",
		Description:   "A VPC is a virtual network that belongs to an account and provides logical isolation from other networks.",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate:       listIsInstanceDisk,
			ParentHydrate: listIsInstance,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIsInstanceDisk,
			KeyColumns: plugin.AllColumns([]string{"id", "instance_id"}),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the instance disk."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The user defined name for this disk."},
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The instance identifier.", Transform: transform.FromField("InstanceId")},
			// Other columns
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time that the disk was created.", Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp)},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "The URL for this instance disk."},
			{Name: "interface_type", Type: proto.ColumnType_STRING, Description: "The disk interface used for attaching the disk."},
			{Name: "resource_type", Type: proto.ColumnType_STRING, Description: "The resource type."},
			{Name: "size", Type: proto.ColumnType_INT, Description: "The size of the disk in GB (gigabytes)."},
			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getAccountId).WithCache(), Transform: transform.FromValue(), Description: "The account ID of this instance disk."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this instance disk."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
		},
	}
}

type instanceDiskInfo = struct {
	vpcv1.InstanceDisk
	InstanceId string
}

//// LIST FUNCTION

func listIsInstanceDisk(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Get instance details
	instanceData := h.Item.(vpcv1.Instance)
	var instanceRegion string
	splitCRN := strings.Split(*instanceData.CRN, ":")
	if len(splitCRN) > 5 {
		instanceRegion = strings.Split(*instanceData.CRN, ":")[5]
	}

	// Return nil, if the config region doesn't contains the instance zone
	if !strings.Contains(instanceRegion, region) {
		return nil, nil
	}

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_instance_disk.listIsInstanceDisk", "connection_error", err)
		return nil, err
	}

	opts := &vpcv1.ListInstanceDisksOptions{
		InstanceID: instanceData.ID,
	}

	result, resp, err := conn.ListInstanceDisksWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_instance_disk.listIsInstanceDisk", "query_error", err, "resp", resp)
		return nil, err
	}
	for _, i := range result.Disks {
		d.StreamListItem(ctx, instanceDiskInfo{i, *instanceData.ID})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIsInstanceDisk(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_instance_disk.getIsInstanceDisk", "connection_error", err)
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetStringValue()
	instanceId := d.KeyColumnQuals["instance_id"].GetStringValue()

	// No inputs
	if id == "" && instanceId == "" {
		return nil, nil
	}

	// Retrieve the get of instance disk for your account
	opts := &vpcv1.GetInstanceDiskOptions{
		ID:         &id,
		InstanceID: &instanceId,
	}

	result, resp, err := conn.GetInstanceDiskWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_instance_disk.getIsInstanceDisk", "query_error", err, "resp", resp)
		if err.Error() == "VSI not found" || err.Error() == "Provided disk ID is not defined to VSI" {
			return nil, nil
		}
		return nil, err
	}
	return instanceDiskInfo{*result, instanceId}, nil
}
