package ibm

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-ibm",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"ibm_account":           tableIbmAccount(ctx),
			"ibm_iam_role":          tableIbmIamRole(ctx),
			"ibm_iam_user":          tableIbmIamUser(ctx),
			"ibm_is_instance":       tableIbmIsInstance(ctx),
			"ibm_is_instance_disk":  tableIbmIsInstanceDisk(ctx),
			"ibm_is_region":         tableIbmIsRegion(ctx),
			"ibm_is_security_group": tableIbmIsSecurityGroup(ctx),
			"ibm_is_subnet":         tableIbmIsSubnet(ctx),
			"ibm_is_vpc":            tableIbmIsVpc(ctx),
			"ibm_kms_key":           tableIbmKmsKey(ctx),
			"ibm_kms_key_ring":      tableIbmKmsKeyRing(ctx),
			"ibm_resource_group":    tableIbmResourceGroup(ctx),
		},
	}
	return p
}
