package hubspot

import (
	"context"
	"io/ioutil"
	"encoding/json"
	"time"
)

// DealsService handles tasks relating to Deals within the
// Hubspot API. 
//
// Hubspot API docs: https://developers.hubspot.com/docs/api/crm/deals
type DealsService service

type DealsPostResponse struct {
	ID         string `json:"id,omitempty"`
	Properties struct {
		Createdate                      time.Time `json:"createdate,omitempty"`
		DaysToClose                     string    `json:"days_to_close,omitempty"`
		DealCurrencyCode                string    `json:"deal_currency_code,omitempty"`
		Dealname                        string    `json:"dealname,omitempty"`
		Dealstage                       string    `json:"dealstage,omitempty"`
		Dealtype                        string    `json:"dealtype,omitempty"`
		HsAllOwnerIds                   string    `json:"hs_all_owner_ids,omitempty"`
		HsClosedAmount                  string    `json:"hs_closed_amount,omitempty"`
		HsClosedAmountInHomeCurrency    string    `json:"hs_closed_amount_in_home_currency,omitempty"`
		HsClosedDealCloseDate           string    `json:"hs_closed_deal_close_date,omitempty"`
		HsClosedDealCreateDate          string    `json:"hs_closed_deal_create_date,omitempty"`
		HsClosedWonCount                string    `json:"hs_closed_won_count,omitempty"`
		HsCreatedate                    time.Time `json:"hs_createdate,omitempty"`
		HsDaysToCloseRaw                string    `json:"hs_days_to_close_raw,omitempty"`
		HsDealStageProbabilityShadow    string    `json:"hs_deal_stage_probability_shadow,omitempty"`
		HsExchangeRate                  string    `json:"hs_exchange_rate,omitempty"`
		HsIsClosed                      string    `json:"hs_is_closed,omitempty"`
		HsIsClosedCount                 string    `json:"hs_is_closed_count,omitempty"`
		HsIsClosedLost                  string    `json:"hs_is_closed_lost,omitempty"`
		HsIsClosedWon                   string    `json:"hs_is_closed_won,omitempty"`
		HsIsDealSplit                   string    `json:"hs_is_deal_split,omitempty"`
		HsIsOpenCount                   string    `json:"hs_is_open_count,omitempty"`
		HsLastmodifieddate              time.Time `json:"hs_lastmodifieddate,omitempty"`
		HsObjectID                      string    `json:"hs_object_id,omitempty"`
		HsObjectSource                  string    `json:"hs_object_source,omitempty"`
		HsObjectSourceID                string    `json:"hs_object_source_id,omitempty"`
		HsObjectSourceLabel             string    `json:"hs_object_source_label,omitempty"`
		HsOpenDealCreateDate            string    `json:"hs_open_deal_create_date,omitempty"`
		HsProjectedAmount               string    `json:"hs_projected_amount,omitempty"`
		HsProjectedAmountInHomeCurrency string    `json:"hs_projected_amount_in_home_currency,omitempty"`
		HsUserIdsOfAllOwners            string    `json:"hs_user_ids_of_all_owners,omitempty"`
		HubspotOwnerAssigneddate        time.Time `json:"hubspot_owner_assigneddate,omitempty"`
		HubspotOwnerID                  string    `json:"hubspot_owner_id,omitempty"`
		Pipeline                        string    `json:"pipeline,omitempty"`
	} `json:"properties,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Archived  bool      `json:"archived,omitempty"`
}


// An API call object for Hubspot, can be used to create Deals, Companies and Contacts.
//
// If `Associations` is left blank, it will not be passed through to HubSpot.
type PostDeals struct {
	Properties   			*PostDealsProperties 	`json:"properties,omitempty"`
	Associations 			*[]Associations 	`json:"associations,omitempty"`
}

// All the relevant properties that can be set on a MeVitae Deal object.
type PostDealsProperties struct {
	CloseDate 				string 			`json:"closedate,omitempty"`
	CreateDate 				string 			`json:"createdate,omitempty"`
	DateOfEnquiry 			string 			`json:"date_of_enquiry,omitempty"`
	DaysToClose 			string 			`json:"days_to_close,omitempty"`
	DealCurrencyCode 		string 			`json:"deal_currency_code,omitempty"`
	DealName 				string 			`json:"dealname,omitempty"`
	DealStage 				string 			`json:"dealstage,omitempty"`
	DealType 				string 			`json:"dealtype,omitempty"`
	HubspotOwnerId 			string 			`json:"hubspot_owner_id,omitempty"`
	Pipeline 				string 			`json:"pipeline,omitempty"`
}

// Assocations to link one HubSpotObj to another.
type Associations struct {
	To    					*To    				`json:"to"`
	Types 					*[]Types 			`json:"types"`
}
// The record or activity you want to associate with the deal, specified by its unique id value.
type To struct {
	Id 						int 				`json:"id"`
}

//The type of the association between the deal and the record/activity. Include the 
//associationCategoryand associationTypeId. Default association type IDs are 
//listed https://developers.hubspot.com/docs/api/crm/associations#association-type-id-values. 

type Types struct {
	AssociationCategory 	string 			`json:"associationCategory"`
	AssociationTypeId 		int 				`json:"associationTypeId"`
}


// PostDeal creates a deal with the properties that will be defined.

func (s *DealsService) PostDeal(ctx context.Context, mapping *PostDeals) (DealsPostResponse, error) {
	req, err := s.client.NewRequest("POST", "crm/v3/objects/deals", mapping)
	var result DealsPostResponse
	if err != nil {
		return result, err
	}

	m := new(PostDeals)
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body) // response body is []byte

    if err := json.Unmarshal(body, &result); err != nil {  // Parse []byte to the go struct pointer
        return result, err
    }
	return result, nil
}