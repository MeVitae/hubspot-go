package hubspot

import (
	"context"
	"fmt"
)

// DealsService handles tasks relating to Deals within the
// Hubspot API. 
//
// Hubspot API docs: https://developers.hubspot.com/docs/api/crm/deals
type DealsService service

// An API call object for Hubspot, can be used to create Deals, Companies and Contacts.
//
// If `Associations` is left blank, it will not be passed through to HubSpot.
type HubSpotDealsObj struct {
	Properties   			*DealsProperties 	`json:"properties,omitempty"`
	Associations 			*[]Associations 	`json:"associations,omitempty"`
}

// All the relevant properties that can be set on a MeVitae Deal object.
type DealsProperties struct {
	CloseDate 				*string 			`json:"closedate,omitempty"`
	CreateDate 				*string 			`json:"createdate,omitempty"`
	DateOfEnquiry 			*string 			`json:"date_of_enquiry,omitempty"`
	DaysToClose 			*string 			`json:"days_to_close,omitempty"`
	DealCurrencyCode 		*string 			`json:"deal_currency_code,omitempty"`
	DealName 				*string 			`json:"dealname,omitempty"`
	DealStage 				*string 			`json:"dealstage,omitempty"`
	DealType 				*string 			`json:"dealtype,omitempty"`
	HubspotOwnerId 			*string 			`json:"hubspot_owner_id,omitempty"`
	Pipeline 				*string 			`json:"pipeline,omitempty"`
}

// Assocations to link one HubSpotObj to another.
type Associations struct {
	To    					*To    				`json:"to"`
	Types 					*[]Types 			`json:"types"`
}
// The record or activity you want to associate with the deal, specified by its unique id value.
type To struct {
	Id 						*int 				`json:"id"`
}

//The type of the association between the deal and the record/activity. Include the 
//associationCategoryand associationTypeId. Default association type IDs are 
//listed https://developers.hubspot.com/docs/api/crm/associations#association-type-id-values. 

type Types struct {
	AssociationCategory 	*string 			`json:"associationCategory"`
	AssociationTypeId 		*int 				`json:"associationTypeId"`
}


// PostDeal creates a deal with the properties that will be defined.
//
// GitHub API docs: https://docs.github.com/enterprise-server@3.12/rest/enterprise-admin/ldap#update-ldap-mapping-for-a-user
//
//meta:operation PATCH /admin/ldap/users/{username}/mapping
func (s *DealsService) PostDeal(ctx context.Context, user string, mapping *HubSpotDealsObj) (*HubSpotDealsObj, *Response, error) {
	req, err := s.client.NewRequest("PATCH", "crm/v3/objects/deals", mapping)
	if err != nil {
		return nil, nil, err
	}

	m := new(HubSpotDealsObj)
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}