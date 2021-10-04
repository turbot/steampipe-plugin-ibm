package ibm

import (
	"context"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableIbmAccountSettings(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_account_settings",
		Description: "Account setting information of IBM Cloud.",
		List: &plugin.ListConfig{
			Hydrate: getAccountSettings,
		},
		Columns: []*plugin.Column{
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "An unique ID of the account."},
			{Name: "restrict_create_service_id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "restrict_create_platform_api_key", Type: proto.ColumnType_STRING, Transform: transform.FromField("RestrictCreatePlatformApikey"), Description: "Indicates whether creating a platform API key is access controlled, or not."},
			{Name: "allowed_ip_addresses", Type: proto.ColumnType_STRING, Description: "The IP addresses and subnets from which IAM tokens can be created for the account."},
			{Name: "entity_tag", Type: proto.ColumnType_STRING, Description: "Version of the account settings."},
			{Name: "mfa", Type: proto.ColumnType_STRING, Description: "Defines the MFA trait for the account."},
			{Name: "session_expiration_in_seconds", Type: proto.ColumnType_STRING, Description: "Defines the session expiration in seconds for the account."},
			{Name: "session_invalidation_in_seconds", Type: proto.ColumnType_STRING, Description: "Defines the period of time in seconds in which a session will be invalidated due  to inactivity."},
			{Name: "history", Type: proto.ColumnType_JSON, Description: "History of the Account Settings."},
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
