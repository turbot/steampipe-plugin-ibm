package ibm

import (
	"context"

	"github.com/IBM/platform-services-go-sdk/atrackerv1"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmAtrackerTarget(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_atracker_target",
		Description:   "IBM Cloud Activity Tracker Target.",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listAtrackerTargets,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getAtrackerTarget,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The uuid of this route resource."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of this route."},
			{Name: "crn", Type: proto.ColumnType_STRING, Description: "The crn of this route type resource.", Transform: transform.FromField("CRN")},
			// Other columns
			{Name: "target_type", Type: proto.ColumnType_STRING, Description: "The type of this target."},
			{Name: "encrypt_key", Type: proto.ColumnType_STRING, Description: "The encryption key used to encrypt events before ATracker services buffer them on storage."},
			{Name: "cos_endpoint", Type: proto.ColumnType_JSON, Description: "Property values for a Cloud Object Storage Endpoint."},

			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN").Transform(crnToAccountID), Description: "The account ID of this target."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this target."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRN").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
		},
	}
}

//// LIST FUNCTION

func listAtrackerTargets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	conn, err := ActivityTrackerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listAtrackerTargets", "connection_error", err)
		return nil, err
	}

	result, resp, err := conn.ListTargetsWithContext(ctx, &atrackerv1.ListTargetsOptions{})
	if err != nil {
		if types.ToString(err.Error()) == "Not Found" {
			return nil, nil
		}
		plugin.Logger(ctx).Error("listAtrackerTargets", "query_error", err, "resp", resp)
		return nil, err
	}
	if result != nil {
		for _, i := range result.Targets {
			d.StreamListItem(ctx, i)
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAtrackerTarget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	conn, err := ActivityTrackerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getAtrackerTarget", "connection_error", err)
		return nil, err
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	opts := &atrackerv1.GetTargetOptions{
		ID: &id,
	}

	result, resp, err := conn.GetTargetWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("getAtrackerTarget", "query_error", err, "resp", resp)
		return nil, err
	}

	return result, nil
}
