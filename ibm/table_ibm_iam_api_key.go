package ibm

import (
	"context"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableIbmAPIKey(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_api_key",
		Description: "API keys in the IBM Cloud account.",
		List: &plugin.ListConfig{
			Hydrate: listAPIKey,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "description", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "entity_tag", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "crn", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "iam_id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "created_at", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "modified_at", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "api_key", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "history", Type: proto.ColumnType_JSON, Description: ""},
		},
	}
}

func listAPIKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := iamService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_api_key.listAPIKey", "connection_error", err)
		return nil, err
	}

	userConn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_role.listIamRole", "connection_error", err)
		return nil, err
	}

	userInfo, err := fetchUserDetails(userConn, 2)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.listIamUser", "connection_error", err)
		return nil, err
	}

	opts := &iamidentityv1.ListAPIKeysOptions{
		AccountID: &userInfo.userAccount,
		Scope:     core.StringPtr("account"),
	}

	result, resp, err := conn.ListAPIKeys(opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_api_key.listAPIKey", "query_error", err, "resp", resp)
		return nil, err
	}
	for _, i := range result.Apikeys {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}
