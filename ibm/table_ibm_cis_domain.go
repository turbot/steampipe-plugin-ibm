package ibm

import (
	"context"
	"strings"

	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/networking-go-sdk/globalloadbalancerv1"
	"github.com/IBM/networking-go-sdk/zonessettingsv1"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmCISDomain(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_cis_domain",
		Description: "IBM CIS Domain",
		List: &plugin.ListConfig{
			Hydrate: listCISDomains,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getCISDomain,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The zone id."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The zone name."},

			// Other columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Hydrate: getAccountId, Description: "An unique ID of the account.", Transform: transform.FromValue()},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time that the zone was created."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time that the zone was updated."},
			{Name: "minimum_tls_version", Type: proto.ColumnType_STRING, Hydrate: getTlsMinimumVersion, Description: "The tls version of the zone.", Transform: transform.FromField("Value")},
			{Name: "original_registrar", Type: proto.ColumnType_STRING, Description: "The original registrar of the zone."},
			{Name: "original_dnshost", Type: proto.ColumnType_STRING, Description: "The original dns host of the zone."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The zone status."},
			{Name: "paused", Type: proto.ColumnType_BOOL, Description: "Whether the zone is in paused state."},
			{Name: "web_application_firewall", Type: proto.ColumnType_STRING, Hydrate: getWebApplicationFirewall, Description: "The web application firewall status.", Transform: transform.FromField("Value")},
			{Name: "dns_records", Type: proto.ColumnType_JSON, Hydrate: getDnsRecords, Description: "DNS records for the domain.", Transform: transform.FromValue()},
			{Name: "original_name_servers", Type: proto.ColumnType_JSON, Description: "The original name servers of the zone."},
			{Name: "name_servers", Type: proto.ColumnType_JSON, Description: "The name servers of the zone."},
			{Name: "global_load_balancer", Type: proto.ColumnType_JSON, Hydrate: getGlobalLoadBalancer, Description: "The global load balancer of the zone.", Transform: transform.FromValue()},

			// Standard columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("ID").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
		},
	}
}

//// LIST FUNCTION

func listCISDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	conn, err := cisZoneService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCISDomains", "connection_error", err)
		return nil, err
	}

	result, resp, err := conn.ListZones(&zonesv1.ListZonesOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("listCISDomains", "query_error", err, "resp", resp)
		if strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		}
		return nil, err
	}
	for _, i := range result.Result {
		d.StreamListItem(ctx, i)
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCISDomain(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create service connection
	conn, err := cisZoneService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getCISDomain", "connection_error", err)
		return nil, err
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// No inputs
	if id == "" {
		return nil, nil
	}

	opts := &zonesv1.GetZoneOptions{
		ZoneIdentifier: &id,
	}

	result, resp, err := conn.GetZone(opts)
	if err != nil {
		plugin.Logger(ctx).Error("getCISDomain", "query_error", err, "resp", resp)
		if strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		}
		return nil, err
	}

	return *result.Result, nil
}

func getTlsMinimumVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(zonesv1.ZoneDetails).ID

	// Create service connection
	conn, err := cisZoneSettingService(ctx, d, *id)
	if err != nil {
		plugin.Logger(ctx).Error("getTlsMinimumVersion", "connection_error", err)
		return nil, err
	}

	tls, resp, err := conn.GetMinTlsVersion(&zonessettingsv1.GetMinTlsVersionOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("getTlsMinimumVersion", "query_error", err, "resp", resp)
		if strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		}
		return nil, err
	}

	return tls.Result, nil
}

func getWebApplicationFirewall(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(zonesv1.ZoneDetails).ID
	// Create service connection
	conn, err := cisZoneSettingService(ctx, d, *id)
	if err != nil {
		plugin.Logger(ctx).Error("getWebApplicationFirewall", "connection_error", err)
		return nil, err
	}

	firewall, resp, err := conn.GetWebApplicationFirewall(&zonessettingsv1.GetWebApplicationFirewallOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("getWebApplicationFirewall", "query_error", err, "resp", resp)
		if strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		}
		return nil, err
	}
	return firewall.Result, nil
}

func getGlobalLoadBalancer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(zonesv1.ZoneDetails).ID
	// Create service connection
	conn, err := cisGlobalLoadBalancerService(ctx, d, *id)
	if err != nil {
		plugin.Logger(ctx).Error("cisGlobalLoadBalancerService", "connection_error", err)
		return nil, err
	}

	loadBalancers, resp, err := conn.ListAllLoadBalancers(&globalloadbalancerv1.ListAllLoadBalancersOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("cisGlobalLoadBalancerService", "query_error", err, "resp", resp)
		if strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		}
		return nil, err
	}
	return loadBalancers.Result, nil
}

func getDnsRecords(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(zonesv1.ZoneDetails).ID
	// Create service connection
	conn, err := cisDnsRecordService(ctx, d, *id)
	if err != nil {
		plugin.Logger(ctx).Error("getDnsRecords", "connection_error", err)
		return nil, err
	}

	dnsRecords, resp, err := conn.ListAllDnsRecords(&dnsrecordsv1.ListAllDnsRecordsOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("getDnsRecords", "query_error", err, "resp", resp)
		if strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		}
		return nil, err
	}
	return dnsRecords.Result, nil
}
