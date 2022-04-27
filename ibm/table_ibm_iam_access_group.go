package ibm

import (
	"context"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmIamAccessGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_access_group",
		Description: "Access Groups in the IBM Cloud account.",
		List: &plugin.ListConfig{
			Hydrate: listAccessGroup,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getAccessGroup,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the access group."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the IAM access group."},
			{Name: "is_federated", Type: proto.ColumnType_BOOL, Description: "This is set to true if rules exist for the group.", Transform: transform.FromField("IsFederated")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp the group was created at.", Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp)},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the IAM access group."},
			{Name: "created_by_id", Type: proto.ColumnType_STRING, Description: "The iam_id of the entity that created the group."},
			{Name: "last_modified_at", Type: proto.ColumnType_TIMESTAMP, Description: "Specifies the date and time, the group las modified.", Transform: transform.FromField("LastModifiedAt").Transform(ensureTimestamp)},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "An url to the given group resource."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "ID of the account that this group belongs to.", Hydrate: plugin.HydrateFunc(getAccountId).WithCache(), Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listAccessGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	conn, err := iamAccessGroupService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_access_group.listAccessGroup", "connection_error", err)
		return nil, err
	}

	// Get account details
	getAccountIdCached := plugin.HydrateFunc(getAccountId).WithCache()
	accountID, err := getAccountIdCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	opts := &iamaccessgroupsv2.ListAccessGroupsOptions{
		AccountID: core.StringPtr(accountID.(string)),
	}

	// Retrieve the list of access group for your account.
	maxResult := int64(100)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < maxResult {
			maxResult = *limit
		}
	}
	opts.SetLimit(maxResult)

	result, resp, err := conn.ListAccessGroups(opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_access_group.listAccessGroup", "query_error", err, "resp", resp)
		return nil, err
	}

	for _, i := range result.Groups {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAccessGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	conn, err := iamAccessGroupService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_access_group.getAccessGroup", "connection_error", err)
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	opts := &iamaccessgroupsv2.GetAccessGroupOptions{
		AccessGroupID: &id,
	}

	item, resp, err := conn.GetAccessGroup(opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_access_group.getAccessGroup", "query_error", err, "resp", resp)
		return nil, err
	}

	return item, nil
}
