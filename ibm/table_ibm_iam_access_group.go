package ibm

import (
	"context"

	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

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
			{Name: "is_federated", Type: proto.ColumnType_BOOL, Description: "", Transform: transform.FromField("IsFederated")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "", Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp)},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the IAM access group."},
			{Name: "created_by_id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "ID of the account that this group belongs to."},
			{Name: "last_modified_at", Type: proto.ColumnType_TIMESTAMP, Description: "Specifies the date and time, the group las modified.", Transform: transform.FromField("LastModifiedAt").Transform(ensureTimestamp)},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("Href")},
		},
	}
}

func listAccessGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := iamAccessGroupService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_access_group.listAccessGroup", "connection_error", err)
		return nil, err
	}

	userConn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_access_group.listAccessGroup", "connection_error", err)
		return nil, err
	}

	userInfo, err := fetchUserDetails(userConn, 2)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_access_group.listAccessGroup", "connection_error", err)
		return nil, err
	}

	opts := &iamaccessgroupsv2.ListAccessGroupsOptions{
		AccountID: &userInfo.userAccount,
	}

	result, resp, err := conn.ListAccessGroups(opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_access_group.listAccessGroup", "query_error", err, "resp", resp)
		return nil, err
	}
	for _, i := range result.Groups {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}

func getAccessGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := iamAccessGroupService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_access_group.getAccessGroup", "connection_error", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	id := quals["id"].GetStringValue()

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
