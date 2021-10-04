package ibm

import (
	"context"

	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev2/controllerv2"
	// "github.com/IBM-Cloud/bluemix-go/models"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableIbmCIS(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_cis",
		Description: "TODO",
		List: &plugin.ListConfig{
			Hydrate: listCisInstance,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getCisInstance,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Resource instance name."},
			{Name: "region_id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "resource_group_id", Type: proto.ColumnType_STRING, Description: "The id of the resource group in which the cis instance is present."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "Unique identifier of resource instance."},
			{Name: "crn", Type: proto.ColumnType_STRING, Description: "The location or the environment in which cis instance exists."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("CRN")},
			// {Name: "guid", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "type", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "state", Type: proto.ColumnType_STRING, Description: ""},
			// {Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			// {Name: "deleted_at", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			// {Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			{Name: "restored_at", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			{Name: "target_crn", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "resource_id", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "allow_cleanup", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "resource_keys_url", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "Resource_group_name", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "siblings_url", Type: proto.ColumnType_STRING, Description: "The state of the user. Possible values are PROCESSING, PENDING, ACTIVE, DISABLED_CLASSIC_INFRASTRUCTURE, and VPN_ONLY."},
			// {Name: "dashboard_url", Type: proto.ColumnType_STRING, Description: "The email of the user."},
			{Name: "resource_group_crn", Type: proto.ColumnType_STRING, Description: "The phone number of the user."},
			{Name: "resource_plan_id", Type: proto.ColumnType_STRING, Description: "The alternative phone number of the user."},
			{Name: "resource_bindings_url", Type: proto.ColumnType_STRING, Description: "A link to a photo of the user."},
			{Name: "resource_aliases_url", Type: proto.ColumnType_STRING, Description: "An alphanumeric value identifying the account ID."},
			{Name: "plan_history", Type: proto.ColumnType_JSON, Description: "User settings."},
			{Name: "last_operation", Type: proto.ColumnType_JSON, Description: "User settings."},
			// {Name: "data", Type: proto.ColumnType_JSON, Description: "The unique ID that identifies an active directory user.", Transform: transform.FromValue()}, // For debugging
		},
	}
}

func listCisInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cis.listCisInstance", "connection_error", err)
		return nil, err
	}

	svc, err := controllerv2.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cis.listCisInstance", "connection_error", err)
		return nil, err
	}

	client := svc.ResourceServiceInstanceV2()

	rsInstQuery := controllerv2.ServiceInstanceQuery{}

	instances, err := client.ListInstances(rsInstQuery)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cis.listCisInstance", "query_error", err)
		return nil, err
	}
	for _, i := range instances {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}

func getCisInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cis.getCisInstance", "connection_error", err)
		return nil, err
	}

	svc, err := controllerv2.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cis.getCisInstance", "connection_error", err)
		return nil, err
	}

	client := svc.ResourceServiceInstanceV2()

	quals := d.KeyColumnQuals
	id := quals["id"].GetStringValue()

	item, err := client.GetInstance(id)
	// plugin.Logger(ctx).Info("ibm_cis.getCisInstance", "query_itemr", item)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_cis.getCisInstance", "query_error", err)
		return nil, err
	}

	return item, nil
}
