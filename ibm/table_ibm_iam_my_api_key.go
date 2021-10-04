package ibm

import (
	"context"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableIbmMyAPIKey(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_my_api_key",
		Description: "User's API key in the IBM Cloud account.",
		List: &plugin.ListConfig{
			Hydrate: getAPIKey,
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

func getAPIKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := iamService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_my_api_key.getAPIKey", "connection_error", err)
		return nil, err
	}
	opts := &iamidentityv1.ListAPIKeysOptions{}

	result, resp, err := conn.ListAPIKeys(opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_my_api_key.getAPIKey", "query_error", err, "resp", resp)
		return nil, err
	}
	for _, i := range result.Apikeys {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}
