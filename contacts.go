package hubspot


import (
	"context"
	"io/ioutil"
	"encoding/json"
	"time"
)


type ContactsService service

type ContactPostResponse struct {
	Status        string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
	CorrelationID string `json:"correlationId,omitempty"`
	Category      string `json:"category,omitempty"`
	ID         string `json:"id,omitempty"`
	Properties struct {
		Createdate                            time.Time `json:"createdate,omitempty"`
		Email                                 string    `json:"email,omitempty"`
		Firstname                             string    `json:"firstname,omitempty"`
		HsAllContactVids                      string    `json:"hs_all_contact_vids,omitempty"`
		HsCurrentlyEnrolledInProspectingAgent string    `json:"hs_currently_enrolled_in_prospecting_agent,omitempty"`
		HsEmailDomain                         string    `json:"hs_email_domain,omitempty"`
		HsIsContact                           string    `json:"hs_is_contact,omitempty"`
		HsIsUnworked                          string    `json:"hs_is_unworked,omitempty"`
		HsLifecyclestageLeadDate              time.Time `json:"hs_lifecyclestage_lead_date,omitempty"`
		HsMarketableStatus                    string    `json:"hs_marketable_status,omitempty"`
		HsMarketableUntilRenewal              string    `json:"hs_marketable_until_renewal,omitempty"`
		HsMembershipHasAccessedPrivateContent string    `json:"hs_membership_has_accessed_private_content,omitempty"`
		HsObjectID                            string    `json:"hs_object_id,omitempty"`
		HsObjectSource                        string    `json:"hs_object_source,omitempty"`
		HsObjectSourceID                      string    `json:"hs_object_source_id,omitempty"`
		HsObjectSourceLabel                   string    `json:"hs_object_source_label,omitempty"`
		HsPipeline                            string    `json:"hs_pipeline,omitempty"`
		HsRegisteredMember                    string    `json:"hs_registered_member,omitempty"`
		Lastmodifieddate                      time.Time `json:"lastmodifieddate,omitempty"`
		Lastname                              string    `json:"lastname,omitempty"`
		Lifecyclestage                        string    `json:"lifecyclestage,omitempty"`
	} `json:"properties,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Archived  bool      `json:"archived,omitempty"`
}


type PostContact struct {
	Properties   			*PostContactProperties 	`json:"properties,omitempty"`
}

type PostContactProperties struct {
	Email				string 			`json:"email,omitempty"`
	FirstName			string 			`json:"firstname,omitempty"`
	LastName 			string 			`json:"lastname,omitempty"`
}


func (s *ContactsService) PostContact(ctx context.Context, mapping *PostContact) (ContactPostResponse, error) {
	req, err := s.client.NewRequest("POST", "crm/v3/objects/contacts", mapping)
	var result ContactPostResponse
	if err != nil {
		return result, err
	}


	m := new(PostContact)
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