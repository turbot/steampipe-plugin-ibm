package ibm

import (
	"context"

	"github.com/IBM-Cloud/bluemix-go/api/certificatemanager"
	"github.com/IBM-Cloud/bluemix-go/models"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIbmCertificateManagerCertificate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_certificate_manager_certificate",
		Description: "Retrieve the details of an existing certificate instance resource and lists all the certificates.",
		List: &plugin.ListConfig{
			Hydrate:    listCertificate,
			KeyColumns: plugin.SingleColumn("certificate_manager_instance_id"),
		},
		Get: &plugin.GetConfig{
			Hydrate:    getCertificate,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID of the certificate that is managed in certificate manager."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the certificate."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the certificate."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of a certificate."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "The serial number of a certificate."},
			{Name: "certificate_manager_instance_id", Type: proto.ColumnType_STRING, Description: "The CRN of the certificate manager service instance.", Transform: transform.FromQual("certificate_manager_instance_id")},
			{Name: "algorithm", Type: proto.ColumnType_STRING, Description: "The Algorithm of a certificate."},
			{Name: "begins_on", Type: proto.ColumnType_INT, Description: "The creation date of the certificate.", Transform: transform.FromField("beginsOn").Transform(ensureTimestamp)},
			{Name: "expires_on", Type: proto.ColumnType_INT, Description: "The expiration date of the certificate.", Transform: transform.FromField("beginsOn").Transform(ensureTimestamp)},
			{Name: "domains", Type: proto.ColumnType_JSON, Description: "An array of valid domains for the issued certificate. The first domain is the primary domain, extra domains are secondary domains."},
			{Name: "has_previous", Type: proto.ColumnType_BOOL, Description: "Indicates whether a certificate has a previous version."},
			{Name: "rotate_keys", Type: proto.ColumnType_BOOL, Description: "Rotate keys."},
			{Name: "imported", Type: proto.ColumnType_BOOL, Description: "Indicates whether a certificate has imported or not."},
			{Name: "issuance_info", Type: proto.ColumnType_JSON, Description: "The issuance information of a certificate."},
			{Name: "issuer", Type: proto.ColumnType_STRING, Description: "The issuer of the certificate."},
			{Name: "key_algorithm", Type: proto.ColumnType_STRING, Description: "An alphanumeric value identifying the account ID."},
			{Name: "data", Type: proto.ColumnType_STRING, Description: "The certificate data.", Hydrate: getCertificate, Transform: transform.FromValue()},
			{Name: "data_key_id", Type: proto.ColumnType_STRING, Description: "The data key id.", Hydrate: getCertificate, Transform: transform.FromValue()},
			{Name: "order_policy", Type: proto.ColumnType_JSON, Description: "The order policy of the certificate."},
		},
	}
}

//// LIST FUNCTION

func listCertificate(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

	instanceID := d.KeyColumnQuals["certificate_manager_instance_id"].GetStringValue()

	certificates, err := client.ListCertificates(instanceID)
	plugin.Logger(ctx).Error("ibm_certificate_manager_certificate.listCertificate", "certificates>>>>>>>>", certificates)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_certificate_manager_certificate.listCertificate", "query_error", err)
		return nil, err
	}
	for _, i := range certificates {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCertificate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_certificate_manager_certificate.getCertificate", "connection_error", err)
		return nil, err
	}

	svc, err := certificatemanager.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_certificate_manager_certificate.getCertificate", "connection_error", err)
		return nil, err
	}

	client := svc.Certificate()

	var id string
	if h.Item != nil {
		id = h.Item.(models.CertificateInfo).ID
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	certificate, err := client.GetCertData(id)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_certificate_manager_certificate.getCertificate", "query_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Warn("getCertificate", "certificate", certificate)

	return certificate, nil
}
