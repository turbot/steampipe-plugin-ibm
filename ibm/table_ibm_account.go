package ibm

import (
	"context"
	"time"

	"github.com/IBM-Cloud/bluemix-go/api/account/accountv1"
	"github.com/IBM-Cloud/bluemix-go/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			{Name: "account_id", Type: proto.ColumnType_STRING, Hydrate: getAccountId, Transform: transform.FromValue(), Description: "The ID fof the account."},
		},
	}
}

type AccountInfo struct {
	BillingCountryCode   string
	BluemixSubscriptions []models.BluemixSubscription
	ConfigurationID      *string
	CountryCode          string
	CurrencyCode         string
	CurrentBillingSystem string
	CustomerID           string
	IsIBMer              bool
	Linkages             []models.AccountLinkage
	Name                 string
	OfferTemplate        string
	Onboarded            int
	OrganizationsRegion  []models.OrganizationsRegion
	Origin               string
	Owner                string
	OwnerIAMID           string
	OwnerUniqueID        string
	OwnerUserID          string
	State                string
	SubscriptionID       string
	Tags                 []interface{}
	TeamDirectoryEnabled bool
	TermsAndConditions   models.TermsAndConditions
	Type                 string
	CreatedAt            time.Time
	GUID                 string
	UpdateComments       string
	UpdatedAt            time.Time
	UpdatedBy            string
	URL                  string
}

//// LIST FUNCTION

func listAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_account.listAccount", "connection_error", err)
		return nil, err
	}

	svc, err := accountv1.New(conn)
	if err != nil {
		plugin.Logger(ctx).Error("ibm_account.listAccount", "connection_error", err)
		return nil, err
	}

	client := svc.Accounts()

	accountID, err := getAccountId(ctx, d, h)
	if err != nil {
		return nil, err
	}

	data, err := client.Get(accountID.(string))
	if err != nil {
		plugin.Logger(ctx).Error("ibm_account.listAccount", "query_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, &AccountInfo{
		BillingCountryCode:   data.Entity.BillingCountryCode,
		BluemixSubscriptions: data.Entity.BluemixSubscriptions,
		ConfigurationID:      data.Entity.ConfigurationID,
		CountryCode:          data.Entity.CountryCode,
		CurrencyCode:         data.Entity.CurrencyCode,
		CurrentBillingSystem: data.Entity.CurrentBillingSystem,
		CustomerID:           data.Entity.CustomerID,
		IsIBMer:              data.Entity.IsIBMer,
		Linkages:             data.Entity.Linkages,
		Name:                 data.Entity.Name,
		OfferTemplate:        data.Entity.OfferTemplate,
		Onboarded:            data.Entity.Onboarded,
		OrganizationsRegion:  data.Entity.OrganizationsRegion,
		Origin:               data.Entity.Origin,
		Owner:                data.Entity.Owner,
		OwnerIAMID:           data.Entity.OwnerIAMID,
		OwnerUniqueID:        data.Entity.OwnerUniqueID,
		OwnerUserID:          data.Entity.OwnerUserID,
		State:                data.Entity.State,
		SubscriptionID:       data.Entity.SubscriptionID,
		Tags:                 data.Entity.Tags,
		TeamDirectoryEnabled: data.Entity.TeamDirectoryEnabled,
		TermsAndConditions:   data.Entity.TermsAndConditions,
		Type:                 data.Entity.Type,
		CreatedAt:            data.Metadata.CreatedAt,
		GUID:                 data.Metadata.GUID,
		UpdateComments:       data.Metadata.UpdateComments,
		UpdatedAt:            data.Metadata.UpdatedAt,
		UpdatedBy:            data.Metadata.UpdatedBy,
		URL:                  data.Metadata.URL,
	})

	return nil, nil
}
