package ibm

import (
	"context"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableIbmIamUserPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_user_policy",
		Description: "Access Groups in the IBM Cloud account.",
		List: &plugin.ListConfig{
			Hydrate: listUserPolicy,
		},
		// Get: &plugin.GetConfig{
		// 	Hydrate:    getAccessGroup,
		// 	KeyColumns: plugin.SingleColumn("id"),
		// },
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the IAM user policy."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "type", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "", Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp)},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the IAM access group."},
			{Name: "created_by_id", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "account_id", Type: proto.ColumnType_STRING, Description: "ID of the account that this group belongs to."},
			{Name: "last_modified_at", Type: proto.ColumnType_TIMESTAMP, Description: "Specifies the date and time, the group las modified.", Transform: transform.FromField("LastModifiedAt").Transform(ensureTimestamp)},
			{Name: "last_modified_by_id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "resources", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "subjects", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "roles", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("Href")},
		},
	}
}

func listUserPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := iamUserPolicy(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user_policy.listUserPolicy", "connection_error", err)
		return nil, err
	}

	userConn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user_policy.listUserPolicy", "connection_error", err)
		return nil, err
	}

	userInfo, err := fetchUserDetails(userConn, 2)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user_policy.listUserPolicy", "connection_error", err)
		return nil, err
	}

	opts := &iampolicymanagementv1.ListPoliciesOptions{
		AccountID: &userInfo.userAccount,
		Type:      core.StringPtr("access"),
	}

	result, resp, err := conn.ListPolicies(opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user_policy.listUserPolicy", "query_error", err, "resp", resp)
		return nil, err
	}
	for _, i := range result.Policies {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}

// func getUserPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

// 	conn, err := iamAccessGroupService(ctx, d)
// 	if err != nil {
// 		plugin.Logger(ctx).Error("ibm_iam_user_policy.getAccessGroup", "connection_error", err)
// 		return nil, err
// 	}

// 	quals := d.KeyColumnQuals
// 	id := quals["id"].GetStringValue()

// 	opts := &iamaccessgroupsv2.GetAccessGroupOptions{
// 		AccessGroupID: &id,
// 	}

// 	item, resp, err := conn.GetAccessGroup(opts)
// 	if err != nil {
// 		plugin.Logger(ctx).Error("ibm_iam_user_policy.getAccessGroup", "query_error", err, "resp", resp)
// 		return nil, err
// 	}

// 	return item, nil
// }
