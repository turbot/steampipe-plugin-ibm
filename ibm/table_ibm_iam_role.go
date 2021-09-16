package ibm

import (
	"context"

	"github.com/IBM-Cloud/bluemix-go/api/iampap/iampapv2"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	//"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableIbmIamRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_role",
		Description: "TODO",
		List: &plugin.ListConfig{
			Hydrate: listIamRole,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIamRole,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "crn", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			{Name: "created_by_id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "last_modified_at", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			{Name: "last_modified_by_id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "name", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "service_name", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "An alphanumeric value identifying the account ID."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "description", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "actions", Type: proto.ColumnType_JSON, Description: ""},
		},
	}
}

func listIamRole(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_role.listIamRole", "connection_error", err)
		return nil, err
	}

	svc, err := iampapv2.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_role.listIamRole", "connection_error", err)
		return nil, err
	}

	client := svc.IAMRoles()

	userInfo, err := fetchUserDetails(conn, 2)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_user.listIamUser", "connection_error", err)
		return nil, err
	}

	opts := iampapv2.RoleQuery{AccountID: userInfo.userAccount}
	roles, err := client.ListAll(opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_role.listIamRole", "query_error", err)
		return nil, err
	}
	for _, i := range roles {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}

func getIamRole(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_role.getIamRole", "connection_error", err)
		return nil, err
	}

	svc, err := iampapv2.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_role.getIamRole", "connection_error", err)
		return nil, err
	}

	client := svc.IAMRoles()

	quals := d.KeyColumnQuals
	id := quals["id"].GetStringValue()

	role, etag, err := client.Get(id)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_role.getIamRole", "query_error", err, "etag", etag)
		return nil, err
	}

	plugin.Logger(ctx).Warn("getIamRole", "role", role)

	return role, nil
}
