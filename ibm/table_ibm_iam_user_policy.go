package ibm

import (
	"context"

	"github.com/IBM-Cloud/bluemix-go/api/usermanagement/usermanagementv2"
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmIamUserPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_user_policy",
		Description: "User access policies in the IBM Cloud account.",
		List: &plugin.ListConfig{
			Hydrate:       listUserPolicy,
			ParentHydrate: listIamUser,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the IAM user policy."},
			{Name: "iam_id", Type: proto.ColumnType_STRING, Description: "An alphanumeric value identifying the user's IAM ID."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The policy type."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the policy was created.", Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp)},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the IAM access group."},
			{Name: "created_by_id", Type: proto.ColumnType_STRING, Description: "The iam ID of the entity that created the policy.", Transform: transform.FromField("CreatedByID")},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "The href link back to the policy.", Transform: transform.FromField("Href")},
			{Name: "last_modified_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the policy was last modified.", Transform: transform.FromField("LastModifiedAt").Transform(ensureTimestamp)},
			{Name: "last_modified_by_id", Type: proto.ColumnType_STRING, Description: "The iam ID of the entity that last modified the policy.", Transform: transform.FromField("LastModifiedByID")},
			{Name: "resources", Type: proto.ColumnType_JSON, Description: "The resources associated with a policy."},
			{Name: "subjects", Type: proto.ColumnType_JSON, Description: "The subjects associated with a policy."},
			{Name: "roles", Type: proto.ColumnType_JSON, Description: "A set of role cloud resource names (CRNs) granted by the policy."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "ID of the account that this policy belongs to.", Hydrate: plugin.HydrateFunc(getAccountId).WithCache(), Transform: transform.FromValue()},
		},
	}
}

type userAccessPolicy struct {
	iampolicymanagementv1.Policy
	IamID string
}

//// LIST FUNCTION

func listUserPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	conn, err := iamPolicyManagementService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user_policy.listUserPolicy", "connection_error", err)
		return nil, err
	}

	// Get user details
	userData := h.Item.(usermanagementv2.UserInfo)

	// Get account details
	getAccountIdCached := plugin.HydrateFunc(getAccountId).WithCache()
	accountID, err := getAccountIdCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	opts := &iampolicymanagementv1.ListPoliciesOptions{
		AccountID: core.StringPtr(accountID.(string)),
		Type:      core.StringPtr("access"),
		IamID:     core.StringPtr(userData.IamID),
	}

	result, resp, err := conn.ListPoliciesWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user_policy.listUserPolicy", "query_error", err, "resp", resp)
		return nil, err
	}

	for _, i := range result.Policies {
		d.StreamListItem(ctx, userAccessPolicy{i, userData.IamID})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
