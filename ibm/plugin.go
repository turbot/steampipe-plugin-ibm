package ibm

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-ibm",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"ibm_account":                         tableIbmAccount(ctx),
			"ibm_certificate_manager_certificate": tableIbmCertificateManagerCertificate(ctx),
			"ibm_cis_domain":                      tableIbmCISDomain(ctx),
			"ibm_cos_bucket":                      tableCosBucket(ctx),
			"ibm_iam_access_group":                tableIbmIamAccessGroup(ctx),
			"ibm_iam_access_group_policy":         tableIbmIamAccessGroupPolicy(ctx),
			"ibm_iam_account_settings":            tableIbmAccountSettings(ctx),
			"ibm_iam_api_key":                     tableIbmIamAPIKey(ctx),
			"ibm_iam_my_api_key":                  tableIbmIamMyAPIKey(ctx),
			"ibm_iam_role":                        tableIbmIamRole(ctx),
			"ibm_iam_user":                        tableIbmIamUser(ctx),
			"ibm_iam_user_policy":                 tableIbmIamUserPolicy(ctx),
			"ibm_is_flow_log":                     tableIbmIsFlowLog(ctx),
			"ibm_is_instance":                     tableIbmIsInstance(ctx),
			"ibm_is_instance_disk":                tableIbmIsInstanceDisk(ctx),
			"ibm_is_network_acl":                  tableIbmIsNetworkAcl(ctx),
			"ibm_is_region":                       tableIbmIsRegion(ctx),
			"ibm_is_security_group":               tableIbmIsSecurityGroup(ctx),
			"ibm_is_subnet":                       tableIbmIsSubnet(ctx),
			"ibm_is_volume":                       tableIbmIsVolume(ctx),
			"ibm_is_vpc":                          tableIbmIsVpc(ctx),
			"ibm_kms_key":                         tableIbmKmsKey(ctx),
			"ibm_kms_key_ring":                    tableIbmKmsKeyRing(ctx),
			"ibm_resource_group":                  tableIbmResourceGroup(ctx),
		},
	}
	return p
}
