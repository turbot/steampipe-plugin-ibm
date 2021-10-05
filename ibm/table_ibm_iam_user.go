package ibm

import (
	"context"

	"github.com/IBM-Cloud/bluemix-go/api/usermanagement/usermanagementv2"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableIbmIamUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_user",
		Description: "Users in the IBM Cloud account.",
		List: &plugin.ListConfig{
			Hydrate: listIamUser,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIamUser,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "An alphanumeric value identifying the user profile."},
			{Name: "iam_id", Type: proto.ColumnType_STRING, Description: "An alphanumeric value identifying the user's IAM ID."},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "The user ID used for login."},
			{Name: "realm", Type: proto.ColumnType_STRING, Description: "The realm of the user. The value is either IBMid or SL."},
			{Name: "first_name", Type: proto.ColumnType_STRING, Description: "The first name of the user.", Transform: transform.FromField("Firstname")},
			{Name: "last_name", Type: proto.ColumnType_STRING, Description: "The last name of the user.", Transform: transform.FromField("Lastname")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the user. Possible values are PROCESSING, PENDING, ACTIVE, DISABLED_CLASSIC_INFRASTRUCTURE, and VPN_ONLY."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The email of the user."},
			{Name: "phonenumber", Type: proto.ColumnType_STRING, Description: "The phone number of the user."},
			{Name: "alt_phonenumber", Type: proto.ColumnType_STRING, Description: "The alternative phone number of the user.", Transform: transform.FromField("Altphonenumber")},
			{Name: "photo", Type: proto.ColumnType_STRING, Description: "A link to a photo of the user."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "An alphanumeric value identifying the account ID."},
			{Name: "settings", Type: proto.ColumnType_JSON, Hydrate: getIamUserSettings, Transform: transform.FromValue(), Description: "User settings."},
		},
	}
}

func listIamUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.listIamUser", "connection_error", err)
		return nil, err
	}

	svc, err := usermanagementv2.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.listIamUser", "connection_error", err)
		return nil, err
	}

	client := svc.UserInvite()

	getAccountIdCached := plugin.HydrateFunc(getAccountId).WithCache()
	accountID, err := getAccountIdCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	users, err := client.ListUsers(accountID.(string))
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.listIamUser", "query_error", err)
		return nil, err
	}
	for _, i := range users {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}

func getIamUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.getIamUser", "connection_error", err)
		return nil, err
	}

	svc, err := usermanagementv2.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.getIamUser", "connection_error", err)
		return nil, err
	}

	client := svc.UserInvite()

	getAccountIdCached := plugin.HydrateFunc(getAccountId).WithCache()
	accountID, err := getAccountIdCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	id := quals["id"].GetStringValue()

	item, err := client.GetUserProfile(accountID.(string), id)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.getIamUser", "query_error", err)
		return nil, err
	}

	return item, nil
}

func getIamUserSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(usermanagementv2.UserInfo)

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.getIamUserSettings", "connection_error", err)
		return nil, err
	}

	svc, err := usermanagementv2.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.getIamUserSettings", "connection_error", err)
		return nil, err
	}

	client := svc.UserInvite()

	getAccountIdCached := plugin.HydrateFunc(getAccountId).WithCache()
	accountID, err := getAccountIdCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	item, err := client.GetUserSettings(accountID.(string), user.IamID)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.getIamUserSettings", "query_error", err)
		return nil, err
	}

	return item, nil
}
