package ibm

import (
	"context"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmIsInstance(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_is_instance",
		Description:   "A VPC is a virtual network that belongs to an account and provides logical isolation from other networks.",
		GetMatrixItemFunc: BuildRegionList,
		List: &plugin.ListConfig{
			Hydrate: listIsInstance,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getIsInstance,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier for this virtual server instance."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The user-defined name for this virtual server instance (and default system hostname)."},
			{Name: "crn", Type: proto.ColumnType_STRING, Description: "The CRN for this virtual server instance.", Transform: transform.FromField("CRN")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the virtual server instance."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time that the virtual server instance was created.", Transform: transform.FromField("CreatedAt").Transform(ensureTimestamp)},
			// Other columns
			{Name: "bandwidth", Type: proto.ColumnType_INT, Description: "The total bandwidth (in megabits per second) shared across the virtual server instance's network interfaces."},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "The URL for this virtual server instance."},
			{Name: "memory", Type: proto.ColumnType_INT, Description: "The amount of memory, truncated to whole gibibytes."},
			{Name: "boot_volume_attachment", Type: proto.ColumnType_JSON, Description: "Specifies the boot volume attachment."},
			{Name: "disks", Type: proto.ColumnType_JSON, Description: "A collection of the instance's disks."},
			{Name: "floating_ips", Type: proto.ColumnType_JSON, Description: "Floating IPs allow inbound and outbound traffic from the Internet to an instance", Hydrate: getInstanceNetworkInterfaceFloatingIps, Transform: transform.FromValue()},
			{Name: "gpu", Type: proto.ColumnType_JSON, Description: "The virtual server instance GPU configuration."},
			{Name: "image", Type: proto.ColumnType_JSON, Description: "The image the virtual server instance was provisioned from."},
			{Name: "network_interfaces", Type: proto.ColumnType_JSON, Description: "A collection of the virtual server instance's network interfaces, including the primary network interface."},
			{Name: "primary_network_interface", Type: proto.ColumnType_JSON, Description: "Specifies the primary network interface."},
			{Name: "profile", Type: proto.ColumnType_JSON, Description: "The profile for this virtual server instance."},
			{Name: "resource_group", Type: proto.ColumnType_JSON, Description: "The resource group for this instance."},
			{Name: "vcpu", Type: proto.ColumnType_JSON, Description: "The virtual server instance VCPU configuration."},
			{Name: "volume_attachments", Type: proto.ColumnType_JSON, Description: "A collection of the virtual server instance's volume attachments, including the boot volume attachment."},
			{Name: "vpc", Type: proto.ColumnType_JSON, Description: "The VPC this virtual server instance resides in.", Transform: transform.FromField("VPC")},
			{Name: "zone", Type: proto.ColumnType_JSON, Description: "The zone this virtual server instance resides in."},

			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("CRN").Transform(crnToAccountID), Description: "The account ID of this instance."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this instance."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("CRN").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
			{Name: "tags", Type: proto.ColumnType_JSON, Hydrate: getInstanceTags, Transform: transform.FromValue(), Description: resourceInterfaceDescription("tags")},
		},
	}
}

//// LIST FUNCTION

func listIsInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_instance.listIsInstance", "connection_error", err)
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

	opts := &vpcv1.ListInstancesOptions{
		Limit: &maxResult,
	}

	// Additional filters
	if d.EqualsQuals["name"] != nil {
		opts.SetName(d.EqualsQuals["name"].GetStringValue())
	}

	for {
		if start != "" {
			opts.Start = &start
		}
		result, resp, err := conn.ListInstancesWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_is_instance.listIsInstance", "query_error", err, "resp", resp)
			return nil, err
		}
		for _, i := range result.Instances {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
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

func getIsInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)

	// Create service connection
	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_instance.getIsInstance", "connection_error", err)
		return nil, err
	}
	id := d.EqualsQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	// Retrieve the get of vpcs for your account.
	opts := &vpcv1.GetInstanceOptions{
		ID: &id,
	}

	result, resp, err := conn.GetInstanceWithContext(ctx, opts)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_instance.getIsInstance", "query_error", err, "resp", resp)
		if err.Error() == "VSI not found" {
			return nil, nil
		}
		return nil, err
	}
	return *result, nil
}

func getInstanceTags(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	vpc := h.Item.(vpcv1.Instance)
	conn, err := tagService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_instance.getInstanceTags", "connection_error", err)
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
			plugin.Logger(ctx).Error("ibm_is_instance.getInstanceTags", "query_error", err, "resp", resp)
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

func getInstanceNetworkInterfaceFloatingIps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)["region"].(string)
	vpc := h.Item.(vpcv1.Instance)

	conn, err := vpcService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_is_instance.getInstanceNetworkInterfaceFloatingIps", "connection_error", err)
		return nil, err
	}

	networkInterfaces := vpc.NetworkInterfaces
	networkInterfaceFloatingIp := []vpcv1.FloatingIP{}

	for _, networkInterface := range networkInterfaces {
		opts := conn.NewListInstanceNetworkInterfaceFloatingIpsOptions(*vpc.ID, *networkInterface.ID)

		result, resp, err := conn.ListInstanceNetworkInterfaceFloatingIpsWithContext(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("ibm_is_instance.getInstanceNetworkInterfaceFloatingIps", "query_error", err, "resp", resp)
			return nil, err
		}

		networkInterfaceFloatingIp = append(networkInterfaceFloatingIp, result.FloatingIps...)
	}

	return networkInterfaceFloatingIp, nil
}
