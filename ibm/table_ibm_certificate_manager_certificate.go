package ibm

import (
	"context"

	"github.com/IBM-Cloud/bluemix-go/api/certificatemanager"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmCertificateManagerCertificate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:          "ibm_certificate_manager_certificate",
		Description:   "Retrieve the details of an existing certificate instance resource and lists all the certificates.",
		GetMatrixItem: BuildServiceInstanceList,
		List: &plugin.ListConfig{
			Hydrate: listCertificate,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "certificate_manager_instance_id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the certificate that is managed in certificate manager."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the certificate."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the certificate."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of a certificate."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "The serial number of a certificate."},
			{Name: "certificate_manager_instance_id", Type: proto.ColumnType_STRING, Description: "The CRN of the certificate manager service instance.", Transform: transform.From(getServiceInstanceCRN)},
			{Name: "algorithm", Type: proto.ColumnType_STRING, Description: "The Algorithm of a certificate."},
			{Name: "auto_renew_enabled", Type: proto.ColumnType_BOOL, Description: "The automatic renewal status of the certificate.", Transform: transform.FromField("OrderPolicy.AutoRenewEnabled"), Default: false},
			{Name: "begins_on", Type: proto.ColumnType_TIMESTAMP, Description: "The creation date of the certificate.", Transform: transform.FromField("BeginsOn").Transform(transform.UnixMsToTimestamp)},
			{Name: "expires_on", Type: proto.ColumnType_TIMESTAMP, Description: "The expiration date of the certificate.", Transform: transform.FromField("ExpiresOn").Transform(transform.UnixMsToTimestamp)},
			{Name: "domains", Type: proto.ColumnType_JSON, Description: "An array of valid domains for the issued certificate. The first domain is the primary domain, extra domains are secondary domains."},
			{Name: "has_previous", Type: proto.ColumnType_BOOL, Description: "Indicates whether a certificate has a previous version."},
			{Name: "rotate_keys", Type: proto.ColumnType_BOOL, Description: "Rotate keys."},
			{Name: "imported", Type: proto.ColumnType_BOOL, Description: "Indicates whether a certificate has imported or not."},
			{Name: "issuance_info", Type: proto.ColumnType_JSON, Description: "The issuance information of a certificate."},
			{Name: "issuer", Type: proto.ColumnType_STRING, Description: "The issuer of the certificate."},
			{Name: "key_algorithm", Type: proto.ColumnType_STRING, Description: "An alphanumeric value identifying the account ID."},
			{Name: "order_policy_name", Type: proto.ColumnType_STRING, Description: "The order policy name of the certificate.", Transform: transform.FromField("OrderPolicy.Name")},
			// Standard columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID").Transform(crnToAccountID), Description: "The account ID of this certificate."},
			{Name: "region", Type: proto.ColumnType_STRING, Transform: transform.From(getRegion), Description: "The region of this certificate."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: resourceInterfaceDescription("title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("ID").Transform(ensureStringArray), Description: resourceInterfaceDescription("akas")},
		},
	}
}

//// LIST FUNCTION

func listCertificate(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listCertificate")

	instanceCRN := plugin.GetMatrixItem(ctx)["instance_crn"].(string)
	serviceType := plugin.GetMatrixItem(ctx)["service_type"].(string)

	// Invalid service type
	if serviceType != "cloudcerts" {
		return nil, nil
	}

	// Return if specified instanceID not matched
	if d.KeyColumnQuals["certificate_manager_instance_id"] != nil && d.KeyColumnQuals["certificate_manager_instance_id"].GetStringValue() != instanceCRN {
		return nil, nil
	}

	// Create service connection
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_certificate_manager_certificate.listCertificate", "connection_error", err)
		return nil, err
	}

	svc, err := certificatemanager.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_certificate_manager_certificate.listCertificate", "connection_error", err)
		return nil, err
	}

	client := svc.Certificate()

	certificates, err := client.ListCertificates(instanceCRN)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_certificate_manager_certificate.listCertificate", "query_error", err)
		return nil, err
	}
	for _, i := range certificates {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
