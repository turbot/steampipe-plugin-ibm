package ibm

import (
	"context"

	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

func tableIbmAccount(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ibm_account",
		Description: "IBM Cloud Account.",
		List: &plugin.ListConfig{
			Hydrate: listAccount,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Specifies the name of the account."},
			{Name: "guid", Type: proto.ColumnType_STRING, Description: "An unique ID of the account."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the account."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The current state of the account."},
			{Name: "country_code", Type: proto.ColumnType_STRING, Description: "Specifies the country code."},
			{Name: "currency_code", Type: proto.ColumnType_STRING, Description: "Specifies the currency type."},
			{Name: "customer_id", Type: proto.ColumnType_STRING, Description: "The customer ID of the account."},
			{Name: "owner_guid", Type: proto.ColumnType_STRING, Description: "An unique Id of the account owner."},
			{Name: "owner_unique_id", Type: proto.ColumnType_STRING, Description: "An unique identifier of the account owner."},
			{Name: "owner_user_id", Type: proto.ColumnType_STRING, Description: "The owner user ID used for login."},
			{Name: "organizations", Type: proto.ColumnType_JSON, Description: "A list of organizations the account is associated."},
			{Name: "members", Type: proto.ColumnType_JSON, Description: "A list of members associated with this account."},
		},
	}
}

//// LIST FUNCTION

func listAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_account.listAccount", "connection_error", err)
		return nil, err
	}

	svc, err := accountv2.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_account.listAccount", "connection_error", err)
		return nil, err
	}

	client := svc.Accounts()

	getAccountIdCached := plugin.HydrateFunc(getAccountId).WithCache()
	accountID, err := getAccountIdCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	data, err := client.Get(accountID.(string))
	if err != nil {
		plugin.Logger(ctx).Error("ibm_account.listAccount", "query_error", err)
		return nil, err
	}
	d.StreamListItem(ctx, *data)

	return nil, nil
}
