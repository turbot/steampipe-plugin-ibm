package ibm

import (
	"context"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmIsFlowLog(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_is_flow_log",
		Description:   "Flow Logs for VPC enable the collection, storage, and presentation of information about the Internet Protocol (IP) traffic going to and from network interfaces within your Virtual Private Cloud (VPC).",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listIsFlowLogs,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIsFlowLog,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{

			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for this flow log collector."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The unique user-defined name for this flow log collector."},
			{Name: "crn", Type: proto.ColumnType_STRING, Description: "The CRN for this flow log collector.", Transform: transform.FromField("CRN")},
			{Name: "lifecycle_state", Type: proto.ColumnType_STRING, Description: "The lifecycle state of the flow log collector."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time that the flow log collector was created.", Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp)},

			// Other columns
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Indicates whether this collector is active."},
			{Name: "auto_delete", Type: proto.ColumnType_BOOL, Description: "If set to `true`, this flow log collector will be automatically deleted when the target is deleted."},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "The URL for this flow log collector."},
			{Name: "resource_group", Type: proto.ColumnType_JSON, Description: "The resource group for this flow log collector."},
			{Name: "storage_bucket", Type: proto.ColumnType_JSON, Description: "The Cloud Object Storage bucket where the collected flows are logged."},
			{Name: "target", Type: proto.ColumnType_JSON, Description: "The target this collector is collecting flow logs for."},
			{Name: "vpc", Type: proto.ColumnType_JSON, Description: "The VPC this flow log collector is associated with.", Transform: transform.FromField("VPC")},

			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN").Transform(crnToAccountID), Description: "The account ID of this flow log collector."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this flow log collector."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRN").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
			{Name: "tags", Type: proto.ColumnType_JSON, Hydrate: getFlowLogTags, Transform: transform.FromValue(), Description: resourceInterfaceDescription("tags")},
		},
	}
}

//// LIST FUNCTION

func listIsFlowLogs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_flow log.listIsFlowLog", "connection_error", err)
		return nil, err
	}

	// Retrieve the list of flow log collectors for your account.
	maxResult := int64(100)
	start := ""

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < maxResult {
			maxResult = *limit
		}
	}

	opts := &vpcv1.ListFlowLogCollectorsOptions{
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
		result, resp, err := conn.ListFlowLogCollectorsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_is_flow_log.listIsFlowLog", "query_error", err, "resp", resp)
			return nil, err
		}
		for _, i := range result.FlowLogCollectors {
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

func getIsFlowLog(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_flow_log.getIsFlowLog", "connection_error", err)
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	opts := &vpcv1.GetFlowLogCollectorOptions{
		ID: &id,
	}

	result, resp, err := conn.GetFlowLogCollectorWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_flow_log.getIsFlowLog", "query_error", err, "resp", resp)
		if err.Error() == "Flow log collector not found" {
			return nil, nil
		}
		return nil, err
	}
	return *result, nil
}

func getFlowLogTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vpc := h.Item.(vpcv1.FlowLogCollector)
	conn, err := tagService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_flow_log.getFlowLogTags", "connection_error", err)
		return nil, err
	}

	opts := conn.NewListTagsOptions()
	opts.SetLimit(100)
	opts.SetProviders([]string{"ghost"})
	opts.SetOrderByName("asc")
	opts.SetAttachedTo(*vpc.CRN)
	opts.SetOffset(0)

	tags := []string{}

	for {
		result, resp, err := conn.ListTagsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_is_flow_log.getFlowLogTags", "query_error", err, "resp", resp)
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
