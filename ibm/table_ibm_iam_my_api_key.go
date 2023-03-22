package ibm

import (
	"context"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableIbmIamMyAPIKey(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_iam_my_api_key",
		Description: "User's API key in the IBM Cloud account.",
		List: &plugin.ListConfig{
			Hydrate: listMyAPIKey,
		},
		Columns: iamAPIKeyColumns(),
	}
}

//// LIST FUNCTION

func listMyAPIKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := iamService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_my_api_key.listMyAPIKey", "connection_error", err)
		return nil, err
	}
	opts := &iamidentityv1.ListAPIKeysOptions{}

	result, resp, err := conn.ListAPIKeys(opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_iam_my_api_key.listMyAPIKey", "query_error", err, "resp", resp)
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
