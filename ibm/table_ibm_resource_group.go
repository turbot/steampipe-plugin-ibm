package ibm

import (
	"context"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmResourceGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_resource_group",
		Description: "A resource group is a way for you to organize your account resources in customizable groupings",
		List: &plugin.ListConfig{
			Hydrate: listResourceGroup,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:      "is_default",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "An alpha-numeric value identifying the resource group."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The human-readable name of the resource group."},
			{Name: "crn", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN"), Description: "The full CRN (cloud resource name) associated with the resource group."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the resource group."},
			{Name: "is_default", Type: proto.ColumnType_BOOL, Description: "Indicates whether this resource group is default of the account or not.", Transform: transform.FromField("Default")},
			// Other columns
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp), Description: "The date when the resource group was initially created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").Transform(ensureTimestamp), Description: "The date when the resource group was last updated."},
			{Name: "payment_methods_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("PaymentMethodsURL"), Description: "The URL to access the payment methods details that associated with the resource group."},
			{Name: "quota_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("QuotaID"), Description: "An alpha-numeric value identifying the quota ID associated with the resource group."},
			{Name: "quota_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("QuotaURL"), Description: "The URL to access the quota details that associated with the resource group."},
			{Name: "teams_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("TeamsURL"), Description: "The URL to access the team details that associated with the resource group."},
			{Name: "resource_linkages", Type: proto.ColumnType_STRING, Transform: transform.FromField("ResourceLinkages"), Description: "An array of the resources that linked to the resource group."},

			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("AccountID"), Description: "An alpha-numeric value identifying the account ID."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRN").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
			{Name: "tags", Type: proto.ColumnType_JSON, Hydrate: getResourceGroupTags, Transform: transform.FromValue(), Description: resourceInterfaceDescription("tags")},
		},
	}
}

//// LIST FUNCTION

func listResourceGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	conn, err := resourceManagerService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_resource_group.listResourceGroup", "connection_error", err)
		return nil, err
	}

	// Get current account
	getAccountIdCached := plugin.HydrateFunc(getAccountId).WithCache()
	accountID, err := getAccountIdCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	opts := &resourcemanagerv2.ListResourceGroupsOptions{
		AccountID: core.StringPtr(accountID.(string)),
	}

	// Additional filters
	if d.KeyColumnQuals["name"] != nil {
		opts.SetName(d.KeyColumnQuals["name"].GetStringValue())
	}

	if d.KeyColumnQuals["is_default"] != nil {
		opts.SetDefault(d.KeyColumnQuals["is_default"].GetBoolValue())
	}

	// Non-Equals Qual Map handling
	if d.Quals["is_default"] != nil {
		for _, q := range d.Quals["is_default"].Quals {
			value := q.Value.GetBoolValue()
			if q.Operator == "<>" {
				opts.SetDefault(false)
				if !value {
					opts.SetDefault(true)
				}
			}
		}
	}

	result, resp, err := conn.ListResourceGroupsWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_resource_group.listResourceGroup", "query_error", err, "resp", resp)
		return nil, err
	}
	for _, i := range result.Resources {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getResourceGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(resourcemanagerv2.ResourceGroup)
	conn, err := tagService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_resource_group.getResourceGroupTags", "connection_error", err)
		return nil, err
	}

	opts := conn.NewListTagsOptions()
	opts.SetLimit(100)
	opts.SetProviders([]string{"ghost"})
	opts.SetOrderByName("asc")
	opts.SetAttachedTo(*data.CRN)
	opts.SetOffset(0)

	tags := []string{}

	for {
		result, resp, err := conn.ListTagsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_resource_group.getResourceGroupTags", "query_error", err, "resp", resp)
			return nil, err
		}
		for _, i := range result.Items {
			tags = append(tags, *i.Name)
		}
		length := int64(len(tags))
		if length >= *result.TotalCount {
			break
		}
		opts.SetOffset(length)
	}

	return tags, nil
}
