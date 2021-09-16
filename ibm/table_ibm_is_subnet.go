package ibm

import (
	"context"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableIbmIsSubnet(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_is_subnet",
		Description:   "Subnets are contiguous ranges of IP addresses specified in CIDR block notation. Each subnet is within a particular zone and cannot span multiple zones or regions.",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listIsSubnet,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIsSubnet,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for this subnet."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The unique user-defined name for this subnet."},
			// Other columns
			{Name: "available_ipv4_address_count", Type: proto.ColumnType_INT, Description: "The number of IPv4 addresses in this subnet that are not in-use, and have not been reserved by the user or the provider."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp), Description: "The date and time that the subnet was created."},
			{Name: "crn", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN"), Description: "The CRN for this subnet."},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "The URL for this subnet."},
			{Name: "ip_version", Type: proto.ColumnType_STRING, Description: "The IP version(s) supported by this subnet."},
			{Name: "ipv4_cidr_block", Type: proto.ColumnType_CIDR, Transform: transform.FromField("Ipv4CIDRBlock"), Description: "The IPv4 range of the subnet, expressed in CIDR format."},
			{Name: "network_acl", Type: proto.ColumnType_JSON, Description: "The network ACL for this subnet."},
			{Name: "public_gateway", Type: proto.ColumnType_JSON, Description: "The public gateway to handle internet bound traffic for this subnet."},
			{Name: "resource_group", Type: proto.ColumnType_JSON, Description: "The resource group for this subnet."},
			{Name: "routing_table", Type: proto.ColumnType_JSON, Description: "The routing table for this subnet."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of this subnet."},
			{Name: "total_ipv4_address_count", Type: proto.ColumnType_INT, Description: "The total number of IPv4 addresses in this subnet."},
			{Name: "vpc", Type: proto.ColumnType_JSON, Transform: transform.FromField("VPC"), Description: "The VPC this subnet is a part of."},
			{Name: "zone", Type: proto.ColumnType_JSON, Description: "The zone this subnet resides in."},
			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN").Transform(crnToAccountID), Description: "The account ID of this subnet."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this subnet."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRN").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
			{Name: "tags", Type: proto.ColumnType_JSON, Hydrate: getSubnetTags, Transform: transform.FromValue(), Description: resourceInterfaceDescription("tags")},
		},
	}
}

func listIsSubnet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := vpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_subnet.listIsSubnet", "connection_error", err)
		return nil, err
	}

	// Retrieve the list of vpcs for your account.
	limit := int64(100)
	start := ""
	opts := &vpcv1.ListSubnetsOptions{
		Limit: &limit,
	}

	for {
		if start != "" {
			opts.Start = &start
		}
		result, resp, err := conn.ListSubnetsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_is_subnet.listIsSubnet", "query_error", err, "resp", resp)
			return nil, err
		}
		for _, i := range result.Subnets {
			d.StreamListItem(ctx, i)
		}
		start = GetNext(result.Next)
		if start == "" {
			break
		}
	}

	return nil, nil
}

func getIsSubnet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := vpcService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_subnet.getIsSubnet", "connection_error", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	id := quals["id"].GetStringValue()

	// Retrieve the get of vpcs for your account.
	opts := &vpcv1.GetSubnetOptions{
		ID: &id,
	}

	result, resp, err := conn.GetSubnetWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_subnet.getIsSubnet", "query_error", err, "resp", resp)
		if err.Error() == "Subnet not found" {
			return nil, nil
		}
		return nil, err
	}
	return *result, nil
}

func getSubnetTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(vpcv1.Subnet)
	conn, err := tagService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_subnet.getSubnetTags", "connection_error", err)
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
			plugin.Logger(ctx).Error("ibm_is_subnet.getIsSubnet", "query_error", err, "resp", resp)
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
