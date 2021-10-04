package ibm

import (
	"context"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableIbmAccountSettings(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_account_settings",
		Description: "Account setting information of IBM Cloud.",
		List: &plugin.ListConfig{
			Hydrate:    getAccountSettings,
			// KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "restrict_create_service_id", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "entity_tag", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "crn", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "id", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "iam_id", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "account_id", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "created_at", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "modified_at", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "api_key", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "history", Type: proto.ColumnType_JSON, Description: ""},
		},
	}
}

func getAccountSettings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

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
	plugin.Logger(ctx).Warn("&userInfo.userAccount", "role>>>>>>>>>>>>", &userInfo.userAccount)
	opts := &iamidentityv1.GetAccountSettingsOptions{
		AccountID: &userInfo.userAccount,
	}

	accountSettings, resp, err := conn.GetAccountSettings(opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_api_key.listAPIKey", "query_error", err, "resp", resp)
		return nil, err
	}
	plugin.Logger(ctx).Warn("getAccountSettings", "role", *accountSettings.AccountID)

	d.StreamListItem(ctx, accountSettings)

	return nil, nil
}
