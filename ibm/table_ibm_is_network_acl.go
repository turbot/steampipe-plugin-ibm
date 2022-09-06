package ibm

import (
	"context"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmIsNetworkAcl(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_is_network_acl",
		Description:   "Network ACL is an optional layer of security for your VPC that acts as a firewall for controlling traffic in and out of one or more subnets.",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listIsNetworkAcl,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIsNetworkAcl,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for this network ACL"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The user-defined name for this network ACL."},
			// Other columns
			{Name: "crn", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN"), Description: "The CRN for this network ACL."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp), Description: "The date and time that the network ACL was created."},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "The URL for this network ACL."},
			{Name: "resource_group", Type: proto.ColumnType_JSON, Description: "The resource group for this network ACL."},
			{Name: "rules", Type: proto.ColumnType_JSON, Description: "The ordered rules for this network ACL. If no rules exist, all traffic will be denied."},
			{Name: "subnets", Type: proto.ColumnType_JSON, Description: "The subnets to which this network ACL is attached."},
			{Name: "vpc", Type: proto.ColumnType_JSON, Transform: transform.FromField("VPC"), Description: "he VPC this network ACL is a part of."},
			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN").Transform(crnToAccountID), Description: "The account ID of this subnet."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this subnet."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRN").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
			{Name: "tags", Type: proto.ColumnType_JSON, Hydrate: getNetworkAclTags, Transform: transform.FromValue(), Description: resourceInterfaceDescription("tags")},
		},
	}
}

//// LIST FUNCTION

func listIsNetworkAcl(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_network_acl.listIsNetworkAcl", "connection_error", err)
		return nil, err
	}

	// Retrieve the list of network acls for your account.
	maxResult := int64(100)
	start := ""

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < maxResult {
			maxResult = *limit
		}
	}

	opts := &vpcv1.ListNetworkAclsOptions{
		Limit: &maxResult,
	}

	for {
		if start != "" {
			opts.Start = &start
		}
		result, resp, err := conn.ListNetworkAclsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_is_network_acl.listIsNetworkAcl", "query_error", err, "resp", resp)
			return nil, err
		}
		for _, i := range result.NetworkAcls {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		start = GetNext(result.Next)
		if start == "" {
			break
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIsNetworkAcl(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_network_acl.getIsNetworkAcl", "connection_error", err)
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	// Retrieve the get of vpcs for your account.
	opts := &vpcv1.GetNetworkACLOptions{
		ID: &id,
	}

	result, resp, err := conn.GetNetworkACLWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_network_acl.getIsNetworkAcl", "query_error", err, "resp", resp)
		if err.Error() == "Subnet not found" {
			return nil, nil
		}
		return nil, err
	}
	return *result, nil
}

func getNetworkAclTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(vpcv1.NetworkACL)

	// Create service connection
	conn, err := tagService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_network_acl.getNetworkAclTags", "connection_error", err)
		return nil, err
	}

	opts := conn.NewListTagsOptions()
	opts.SetLimit(100)
	opts.SetProviders([]string{"ghost"})
	opts.SetOrderByName("asc")
	opts.SetAttachedTo(*item.CRN)
	opts.SetOffset(0)

	tags := []string{}

	for {
		result, resp, err := conn.ListTagsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_is_network_acl.getNetworkAclTags", "query_error", err, "resp", resp)
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
