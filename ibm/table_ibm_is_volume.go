package ibm

import (
	"context"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmIsVolume(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_is_volume",
		Description:   "VPC block storage volume.",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listIsVolumes,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIsVolume,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for this volume."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The user-defined name for this volume."},
			{Name: "crn", Type: proto.ColumnType_STRING, Description: "The CRN for this volume.", Transform: transform.FromField("CRN")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the volume."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time that the volume was created.", Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp)},
			// Other columns
			{Name: "capacity", Type: proto.ColumnType_INT, Description: "The capacity of the volume in gigabytes."},
			{Name: "encryption", Type: proto.ColumnType_STRING, Description: "The type of encryption used on the volume."},
			{Name: "encryption_key", Type: proto.ColumnType_STRING, Description: "A reference to the root key used to wrap the data encryption key for the volume. This property will be present for volumes with an `encryption` type of `user_managed`."},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "The URL for this volume."},
			{Name: "iops", Type: proto.ColumnType_INT, Description: "The bandwidth for the volume."},
			{Name: "profile", Type: proto.ColumnType_JSON, Description: "The profile for this volume."},
			{Name: "resource_group", Type: proto.ColumnType_JSON, Description: "The resource group for this volume."},
			{Name: "status_reasons", Type: proto.ColumnType_JSON, Description: "The enumerated reason code values for this property will expand in the future."},
			{Name: "volume_attachments", Type: proto.ColumnType_JSON, Description: "The collection of volume attachments attaching instances to the volume.."},
			{Name: "zone", Type: proto.ColumnType_JSON, Description: "The zone this volume resides in."},

			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN").Transform(crnToAccountID), Description: "The account ID of this volume."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this volume."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRN").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
			{Name: "tags", Type: proto.ColumnType_JSON, Hydrate: getVolumeTags, Transform: transform.FromValue(), Description: resourceInterfaceDescription("tags")},
		},
	}
}

//// LIST FUNCTION

func listIsVolumes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("listIsVolumes", "connection_error", err)
		return nil, err
	}

	// Retrieve the list of VPC block storage volumes for your account.
	maxResult := int64(100)
	start := ""

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < maxResult {
			maxResult = *limit
		}
	}

	opts := &vpcv1.ListVolumesOptions{
		Limit: &maxResult,
	}

	// Additional filters
	if d.KeyColumnQuals["name"] != nil {
		opts.SetName(d.KeyColumnQuals["name"].GetStringValue())
	}

	for {
		if start != "" {
			opts.Start = &start
		}
		result, resp, err := conn.ListVolumesWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("listIsVolumes", "query_error", err, "resp", resp)
			return nil, err
		}
		for _, i := range result.Volumes {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		start = GetNext(result.Next)
		if start == "" {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIsVolume(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("getIsVolume", "connection_error", err)
		return nil, err
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	// Retrieve the vpc block storage volume for your account.
	opts := &vpcv1.GetVolumeOptions{
		ID: &id,
	}

	result, resp, err := conn.GetVolumeWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("getIsVolume", "query_error", err, "resp", resp)
		if err.Error() == "Volume not found" {
			return nil, nil
		}
		return nil, err
	}
	return *result, nil
}

func getVolumeTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	volume := h.Item.(vpcv1.Volume)
	conn, err := tagService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getVolumeTags", "connection_error", err)
		return nil, err
	}

	opts := conn.NewListTagsOptions()
	opts.SetLimit(100)
	opts.SetProviders([]string{"ghost"})
	opts.SetOrderByName("asc")
	opts.SetAttachedTo(*volume.CRN)
	opts.SetOffset(0)

	tags := []string{}

	for {
		result, resp, err := conn.ListTagsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("getVolumeTags", "query_error", err, "resp", resp)
			return nil, err
		}
		for _, i := range result.Items {
			tags = append(tags, *i.Name)
		}
		length := int64(len(tags))
		if length >= *result.TotalCount {
			break
		}
		opts.SetOffset(length)
	}

	return tags, nil
}
