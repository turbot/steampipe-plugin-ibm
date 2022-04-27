package ibm

import (
	"context"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmIsSecurityGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_is_security_group",
		Description:   "Security groups provide a way to apply IP filtering rules to instances in the associated VPC. With security groups, all traffic is denied by default, and rules added to security groups define which traffic the security group permits. Security group rules are stateful such that reverse traffic in response to allowed traffic is automatically permitted.",
		GetMatrixItem: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listIsSecurityGroup,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIsSecurityGroup,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for this security group."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The unique user-defined name for this security group."},
			// Other columns
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp), Description: "The date and time that the security group was created."},
			{Name: "crn", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN"), Description: "The CRN for this security group."},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "The URL for this security group."},
			{Name: "network_interfaces", Type: proto.ColumnType_JSON, Description: "Array of references to network interfaces."},
			{Name: "resource_group", Type: proto.ColumnType_JSON, Description: "The resource group for this security group."},
			{Name: "rules", Type: proto.ColumnType_JSON, Description: "Array of rules for this security group. If no rules exist, all traffic will be denied."},
			{Name: "targets", Type: proto.ColumnType_JSON, Description: "Array of references to targets."},
			{Name: "vpc", Type: proto.ColumnType_JSON, Transform: transform.FromField("VPC"), Description: "The VPC this security group is a part of."},
			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN").Transform(crnToAccountID), Description: "The account ID of this security group."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this security group."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRN").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
			{Name: "tags", Type: proto.ColumnType_JSON, Hydrate: getSecurityGroupTags, Transform: transform.FromValue(), Description: resourceInterfaceDescription("tags")},
		},
	}
}

//// LIST FUNCTION

func listIsSecurityGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_security_group.listIsSecurityGroup", "connection_error", err)
		return nil, err
	}

	// Retrieve the list of vpcs for your account.
	maxResult := int64(100)
	start := ""

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < maxResult {
			maxResult = *limit
		}
	}

	opts := &vpcv1.ListSecurityGroupsOptions{
		Limit: &maxResult,
	}

	for {
		if start != "" {
			opts.Start = &start
		}
		result, resp, err := conn.ListSecurityGroupsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_is_security_group.listIsSecurityGroup", "query_error", err, "resp", resp)
			return nil, err
		}
		for _, i := range result.SecurityGroups {
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

func getIsSecurityGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_security_group.getIsSecurityGroup", "connection_error", err)
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	// Retrieve the get of vpcs for your account.
	opts := &vpcv1.GetSecurityGroupOptions{
		ID: &id,
	}

	result, resp, err := conn.GetSecurityGroupWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_security_group.getIsSecurityGroup", "query_error", err, "resp", resp)
		if err.Error() == "SecurityGroup not found" {
			return nil, nil
		}
		return nil, err
	}
	return *result, nil
}

func getSecurityGroupTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(vpcv1.SecurityGroup)
	conn, err := tagService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_security_group.getSecurityGroupTags", "connection_error", err)
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
			plugin.Logger(ctx).Error("ibm_is_security_group.getIsSecurityGroup", "query_error", err, "resp", resp)
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
