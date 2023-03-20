package ibm

import (
	"context"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func iamAPIKeyColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "name", Type: proto.ColumnType_STRING, Description: "Specifies the name of the API key."},
		{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier of this API Key."},
		{Name: "crn", Type: proto.ColumnType_STRING, Description: "Cloud Resource Name of the API key.", Transform: transform.FromField("CRN")},
		{Name: "iam_id", Type: proto.ColumnType_STRING, Description: "The iam_id that this API key authenticates."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Specifies the date and time, the API key is created.", Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp)},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the API key."},
		{Name: "entity_tag", Type: proto.ColumnType_STRING, Description: "Version of the API Key details object."},
		{Name: "account_id", Type: proto.ColumnType_STRING, Description: "ID of the account that this API key authenticates for."},
		{Name: "modified_at", Type: proto.ColumnType_TIMESTAMP, Description: "Specifies the date and time, the API key las modified.", Transform: transform.FromField("ModifiedAt").Transform(ensureTimestamp)},
		{Name: "api_key", Type: proto.ColumnType_STRING, Description: "The API key value. This property only contains the API key value for the following cases: create an API key, update a service ID API key that stores the API key value as retrievable, or get a service ID API key that stores the API key value as retrievable.", Transform: transform.FromField("Apikey")},
		{Name: "history", Type: proto.ColumnType_JSON, Description: "History of the API key."},
	}
}

//// TABLE DEFINITION

func tableIbmIamAPIKey(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_api_key",
		Description: "API keys in the IBM Cloud account.",
		List: &plugin.ListConfig{
			Hydrate: listAPIKey,
		},
		Columns: iamAPIKeyColumns(),
	}
}

//// LIST FUNCTION

func listAPIKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := iamService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_api_key.listAPIKey", "connection_error", err)
		return nil, err
	}

	// Get account details
	getAccountIdCached := plugin.HydrateFunc(getAccountId).WithCache()
	accountID, err := getAccountIdCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	opts := &iamidentityv1.ListAPIKeysOptions{
		AccountID: core.StringPtr(accountID.(string)),
		Scope:     core.StringPtr("account"),
		Type:      core.StringPtr("user"),
	}

	result, resp, err := conn.ListAPIKeys(opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_api_key.listAPIKey", "query_error", err, "resp", resp)
		return nil, err
	}
	for _, i := range result.Apikeys {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
