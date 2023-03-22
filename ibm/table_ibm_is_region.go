package ibm

import (
	"context"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmIsRegion(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_is_region",
		Description: "IBM Cloud regions and endpoints.",
		List: &plugin.ListConfig{
			Hydrate: listIsRegion,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIsRegion,
			KeyColumns: plugin.SingleColumn("name"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "endpoint", Type: proto.ColumnType_STRING, Description: "The API endpoint for this region."},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "The URL for this region."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The unique user-defined name for this region."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of this region."},
			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getAccountId).WithCache(), Transform: transform.FromValue(), Description: "The account ID of this region."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The region of this region."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			// TODO - should be in crn: format?
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Endpoint").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
			{Name: "tags", Type: proto.ColumnType_JSON, Transform: transform.FromConstant(map[string]string{}), Description: resourceInterfaceDescription("tags")},
		},
	}
}

//// LIST FUNCTION

func listIsRegion(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := GetDefaultIBMRegion(d)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_region.listIsRegion", "connection_error", err)
		return nil, err
	}

	// Retrieve the list of vpcs for your account.
	opts := &vpcv1.ListRegionsOptions{}

	result, resp, err := conn.ListRegionsWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_region.listIsRegion", "query_error", err, "resp", resp)
		return nil, err
	}
	for _, i := range result.Regions {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIsRegion(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := GetDefaultIBMRegion(d)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_region.getIsRegion", "connection_error", err)
		return nil, err
	}
	name := d.EqualsQuals["name"].GetStringValue()

	// No inputs
	if name == "" {
		return nil, nil
	}

	// Retrieve the get of vpcs for your account.
	opts := &vpcv1.GetRegionOptions{
		Name: &name,
	}

	result, resp, err := conn.GetRegionWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_region.getIsRegion", "query_error", err, "resp", resp)
		if err.Error() == "Region not found" {
			return nil, nil
		}
		return nil, err
	}
	return *result, nil
}
