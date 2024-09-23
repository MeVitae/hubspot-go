package hubspot

import (
	"context"
	"io/ioutil"
	"encoding/json"
	"time"
)

type CompaniesService service

type SearchCompanies struct {
	Query        			string   				`json:"query,omitempty"`
	Limit        			int      				`json:"limit,omitempty"`
	After        			string   				`json:"after,omitempty"`
	Sorts        			[]string 				`json:"sorts,omitempty"`
	Properties   			[]string 				`json:"properties,omitempty"`
	FilterGroups 			[]SearchFilterGroups 	`json:"filterGroups,omitempty"`
}
type SearchFilterGroups struct {
	Filters 				[]SearchFilters 		`json:"filters,omitempty"`
}

type SearchFilters struct {
	HighValue    			string   				`json:"highValue,omitempty"`
	PropertyName 			string   				`json:"propertyName,omitempty"`
	Values       			[]string 				`json:"values,omitempty"`
	Value        			string   				`json:"value,omitempty"`
	Operator     			string   				`json:"operator,omitempty"`
}

type SeachCompaniesResponse struct {
	Total   int `json:"total,omitempty"`
	Results []struct {
		ID         string `json:"id,omitempty"`
		Properties struct {
			Createdate         time.Time `json:"createdate,omitempty"`
			Domain             string    `json:"domain,omitempty"`
			HsLastmodifieddate time.Time `json:"hs_lastmodifieddate,omitempty"`
			HsObjectID         string    `json:"hs_object_id,omitempty"`
			Name               string    `json:"name,omitempty"`
		} `json:"properties,omitempty"`
		CreatedAt time.Time `json:"createdAt,omitempty"`
		UpdatedAt time.Time `json:"updatedAt,omitempty"`
		Archived  bool      `json:"archived,omitempty"`
	} `json:"results,omitempty"`
}

type PostCompany struct {
	Properties   			*PostCompanyProperties 	`json:"properties,omitempty"`
}

type PostCompanyProperties struct {
	Domain					string 			`json:"domain,omitempty"`
	Name					string 			`json:"name,omitempty"`
	ATS						string 			`json:"ats,omitempty"`
	CompanySize				string 			`json:"company_size,omitempty"`
	TimeZone				string 			`json:"timezone",omitempty"`
}

type PostCompanyResponse struct {
	ID         string `json:"id,omitempty"`
	Properties struct {
		Ats                 string    `json:"ats,omitempty"`
		CompanySize         string    `json:"company_size,omitempty"`
		Createdate          time.Time `json:"createdate,omitempty"`
		Domain              string    `json:"domain,omitempty"`
		HsLastmodifieddate  time.Time `json:"hs_lastmodifieddate,omitempty"`
		HsObjectID          string    `json:"hs_object_id,omitempty"`
		HsObjectSource      string    `json:"hs_object_source,omitempty"`
		HsObjectSourceID    string    `json:"hs_object_source_id,omitempty"`
		HsObjectSourceLabel string    `json:"hs_object_source_label,omitempty"`
		HsPipeline          string    `json:"hs_pipeline,omitempty"`
		Lifecyclestage      string    `json:"lifecyclestage,omitempty"`
		Name                string    `json:"name,omitempty"`
		Website             string    `json:"website,omitempty"`
	} `json:"properties,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Archived  bool      `json:"archived,omitempty"`
}

func (s *CompaniesService) SearchCompanies(ctx context.Context, mapping *SearchCompanies) (SeachCompaniesResponse, error) {
	req, err := s.client.NewRequest("POST", "crm/v3/objects/companies/search", mapping)
	var result SeachCompaniesResponse
	if err != nil {
		return result, err
	}

	m := new(SearchCompanies)
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

func (s *CompaniesService) PostCompany(ctx context.Context, mapping *PostCompany) (PostCompanyResponse, error) {
	req, err := s.client.NewRequest("POST", "crm/v3/objects/companies", mapping)
	var result PostCompanyResponse
	if err != nil {
		return result, err
	}

	m := new(PostCompany)
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