package ibm

import (
	"context"
	"net/url"
	"reflect"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmIsVpc(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_is_vpc",
		Description:   "A VPC is a virtual network that belongs to an account and provides logical isolation from other networks.",
		GetMatrixItemFunc: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listIsVpc,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:      "classic_access",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIsVpc,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for this VPC."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The unique user-defined name for this VPC."},
			// Other columns
			{Name: "address_prefixes", Type: proto.ColumnType_JSON, Description: "Array of all address pool prefixes for this VPC.", Hydrate: getVpcAddressPrefixes, Transform: transform.FromValue()},
			{Name: "classic_access", Type: proto.ColumnType_BOOL, Description: "Indicates whether this VPC is connected to Classic Infrastructure."},
			{Name: "crn", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN"), Description: "The CRN for this VPC."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp), Description: "The date and time that the VPC was created."},
			{Name: "cse_source_ips", Type: proto.ColumnType_JSON, Description: "Array of CSE source IP addresses for the VPC. The VPC will have one CSE source IP address per zone."},
			{Name: "default_network_acl", Type: proto.ColumnType_JSON, Description: "The default network ACL to use for subnets created in this VPC."},
			{Name: "default_routing_table", Type: proto.ColumnType_JSON, Description: "The default routing table to use for subnets created in this VPC."},
			{Name: "default_security_group", Type: proto.ColumnType_JSON, Description: "The default security group to use for network interfaces created in this VPC."},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "The URL for this VPC."},
			{Name: "resource_group", Type: proto.ColumnType_JSON, Description: "The resource group for this VPC."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of this VPC."},
			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN").Transform(crnToAccountID), Description: "The account ID of this VPC."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this VPC."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRN").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
			{Name: "tags", Type: proto.ColumnType_JSON, Hydrate: getTags, Transform: transform.FromValue(), Description: resourceInterfaceDescription("tags")},
		},
	}
}

func GetNext(next interface{}) string {
	if reflect.ValueOf(next).IsNil() {
		return ""
	}
	u, err := url.Parse(reflect.ValueOf(next).Elem().FieldByName("Href").Elem().String())
	if err != nil {
		return ""
	}
	q := u.Query()
	return q.Get("start")
}

//// LIST FUNCTION

func listIsVpc(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_vpc.listIsVpc", "connection_error", err)
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

	opts := &vpcv1.ListVpcsOptions{
		Limit: &maxResult,
	}

	// Equals Qual Map handling
	if d.KeyColumnQuals["classic_access"] != nil {
		opts.SetClassicAccess(d.KeyColumnQuals["classic_access"].GetBoolValue())
	}

	// Non-Equals Qual Map handling
	if d.Quals["classic_access"] != nil {
		for _, q := range d.Quals["classic_access"].Quals {
			value := q.Value.GetBoolValue()
			if q.Operator == "<>" {
				opts.SetClassicAccess(false)
				if !value {
					opts.SetClassicAccess(true)
				}
			}
		}
	}

	for {
		if start != "" {
			opts.Start = &start
		}
		result, resp, err := conn.ListVpcsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_is_vpc.listIsVpc", "query_error", err, "resp", resp)
			return nil, err
		}
		for _, i := range result.Vpcs {
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

func getIsVpc(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_vpc.getIsVpc", "connection_error", err)
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	// Retrieve the get of vpcs for your account.
	opts := &vpcv1.GetVPCOptions{
		ID: &id,
	}

	result, resp, err := conn.GetVPCWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_vpc.getIsVpc", "query_error", err, "resp", resp)
		if err.Error() == "VPC not found" {
			return nil, nil
		}
		return nil, err
	}
	return *result, nil
}

func getTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vpc := h.Item.(vpcv1.VPC)
	conn, err := tagService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_vpc.getTags", "connection_error", err)
		return nil, err
	}

	opts := conn.NewListTagsOptions()
	opts.SetLimit(100)
	opts.SetProviders([]string{"ghost"})
	opts.SetOrderByName("asc")
	opts.SetAttachedTo(*vpc.CRN)
	opts.SetOffset(0)

	tags := []string{}

	for {
		result, resp, err := conn.ListTagsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_is_vpc.getIsVpc", "query_error", err, "resp", resp)
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

func getVpcAddressPrefixes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)
	vpc := h.Item.(vpcv1.VPC)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_vpc.getVpcAddressPrefixes", "connection_error", err)
		return nil, err
	}

	opts := &vpcv1.ListVPCAddressPrefixesOptions{
		VPCID: vpc.ID,
	}

	result, resp, err := conn.ListVPCAddressPrefixesWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_vpc.getVpcAddressPrefixes", "query_error", err, "resp", resp)
		return nil, err
	}
	return result.AddressPrefixes, nil
}
