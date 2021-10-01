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
			"ibm_iam_role":          tableIbmIamRole(ctx),
			"ibm_iam_user":          tableIbmIamUser(ctx),
			"ibm_is_network_acl":    tableIbmIsNetworkAcl(ctx),
			"ibm_is_region":         tableIbmIsRegion(ctx),
			"ibm_is_security_group": tableIbmIsSecurityGroup(ctx),
			"ibm_is_subnet":         tableIbmIsSubnet(ctx),
			"ibm_is_vpc":            tableIbmIsVpc(ctx),
			//"ibm_resource_group":    tableIbmResourceGroup(ctx),
		},
	}
	return p
}
