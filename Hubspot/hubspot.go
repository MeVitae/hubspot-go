package hubspot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)
const (
	defaultBaseURL    = "https://api.hubspot.com/"

	mediaTypeV3 = "*/*"
)

var errNonNilContext = errors.New("context must be non-nil")

// A Client manages communication with the Hubspot API.
type Client struct {
	client   *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to the public Hubspot API
	// BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the GitHub API.
	Deals			   *DealsService
	Companies          *CompaniesService
	Contacts           *ContactsService
}

type service struct {
	client *Client
}

// NewClient returns a new GitHub API client. If a nil httpClient is
// provided, a new http.Client will be used. To use API methods which require
// authentication, either use Client.WithAuthToken or provide NewClient with
// an http.Client that will perform the authentication for you (such as that
// provided by the golang.org/x/oauth2 library).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient2 := *httpClient
	c := &Client{client: &httpClient2}
	c.initialize()
	return c
}

// RequestOption represents an option that can modify an http.Request.
type RequestOption func(req *http.Request)

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}, opts ...RequestOption) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", mediaTypeV3)

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

// WithAuthToken returns a copy of the client configured to use the provided token for the Authorization header.
func (c *Client) WithAuthToken(token string) *Client {
	c2 := c.copy()
	defer c2.initialize()
	transport := c2.client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	c2.client.Transport = roundTripperFunc(
		func(req *http.Request) (*http.Response, error) {
			req = req.Clone(req.Context())
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			return transport.RoundTrip(req)
		},
	)
	return c2
}

func (c *Client) initialize() {
	if c.client == nil {
		c.client = &http.Client{}
	}
	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse(defaultBaseURL)
	}
	c.common.client = c
	c.Deals = (*DealsService)(&c.common)			   
	c.Companies = (*CompaniesService)(&c.common)	        
	c.Contacts = (*ContactsService)(&c.common)	
	
}

// copy returns a copy of the current client. It must be initialized before use.
func (c *Client) copy() *Client {
	// can't use *c here because that would copy mutexes by value.
	clone := Client{
		client:                  &http.Client{},
		BaseURL:                 c.BaseURL,
	}
	if c.client != nil {
		clone.client.Transport = c.client.Transport
		clone.client.CheckRedirect = c.client.CheckRedirect
		clone.client.Jar = c.client.Jar
		clone.client.Timeout = c.client.Timeout
	}
	return &clone
}

func withContext(ctx context.Context, req *http.Request) *http.Request {
	// No-op because App Engine adds context to a request differently.
	return req
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer interface,
// the raw response body will be written to v, without attempting to first
// decode it. If v is nil, and no error happens, the response is returned as is.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	req = withContext(ctx, req)


	resp, err := c.client.Do(req)
	
	return resp, err
}

// roundTripperFunc creates a RoundTripper (transport)
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}