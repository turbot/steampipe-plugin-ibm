package ibm

import (
	"context"

	"github.com/IBM-Cloud/bluemix-go/api/iampap/iampapv2"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableIbmIamRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_role",
		Description: "An IAM role is an Identity and Access Management (IAM) entity with permissions to make IBM cloud service requests.",
		List: &plugin.ListConfig{
			Hydrate: listIamRole,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIamRole,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The role ID."},
			{Name: "crn", Type: proto.ColumnType_STRING, Description: "The Cloud Resource Name (CRN) that uniquely identifies your cloud resources."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the role was created."},
			{Name: "created_by_id", Type: proto.ColumnType_STRING, Description: "The IAM ID of the entity that created the role."},
			{Name: "last_modified_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the role was last modified."},
			{Name: "last_modified_by_id", Type: proto.ColumnType_STRING, Description: "The IAM ID of the entity that last modified the policy."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the role that is used in the CRN."},
			{Name: "service_name", Type: proto.ColumnType_STRING, Description: "The service name."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "An alphanumeric value identifying the account ID."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the role that is shown in the console."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the role."},
			{Name: "actions", Type: proto.ColumnType_JSON, Description: "The actions of the role."},
		},
	}
}

func listIamRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	// Get account details
	getAccountIdCached := plugin.HydrateFunc(getAccountId).WithCache()
	accountID, err := getAccountIdCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	opts := iampapv2.RoleQuery{AccountID: accountID.(string)}
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

	// No inputs
	if id == "" {
		return nil, nil
	}

	role, etag, err := client.Get(id)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_role.getIamRole", "query_error", err, "etag", etag)
		return nil, err
	}

	return role, nil
}
