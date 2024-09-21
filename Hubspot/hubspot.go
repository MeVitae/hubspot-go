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
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)
const (
	Version = "v1.0.0"

	defaultBaseURL    = "https://api.hubspot.com/"
)

var errNonNilContext = errors.New("context must be non-nil")

// A Client manages communication with the Hubspot API.
type Client struct {
	client   *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to the public Hubspoit API
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
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	req.Header.Set(headerAPIVersion, defaultAPIVersion)

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
	c.clientMu.Lock()
	// can't use *c here because that would copy mutexes by value.
	clone := Client{
		client:                  &http.Client{},
		BaseURL:                 c.BaseURL
	}
	if c.client != nil {
		clone.client.Transport = c.client.Transport
		clone.client.CheckRedirect = c.client.CheckRedirect
		clone.client.Jar = c.client.Jar
		clone.client.Timeout = c.client.Timeout
	}
	return &clone
}

// BareDo sends an API request and lets you handle the api response. If an error
// or API Error occurs, the error will contain more information. Otherwise you
// are supposed to read and close the response's Body. If rate limit is exceeded
// and reset time is in the future, BareDo returns *RateLimitError immediately
// without making a network API call.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	req = withContext(ctx, req)

	rateLimitCategory := GetRateLimitCategory(req.Method, req.URL.Path)

	if bypass := ctx.Value(bypassRateLimitCheck); bypass == nil {
		// If we've hit rate limit, don't make further requests before Reset time.
		if err := c.checkRateLimitBeforeDo(req, rateLimitCategory); err != nil {
			return &Response{
				Response: err.Response,
				Rate:     err.Rate,
			}, err
		}
		// If we've hit a secondary rate limit, don't make further requests before Retry After.
		if err := c.checkSecondaryRateLimitBeforeDo(req); err != nil {
			return &Response{
				Response: err.Response,
			}, err
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(url).String()
				return nil, e
			}
		}

		return nil, err
	}

	response := newResponse(resp)

	// Don't update the rate limits if this was a cached response.
	// X-From-Cache is set by https://github.com/gregjones/httpcache
	if response.Header.Get("X-From-Cache") == "" {
		c.rateMu.Lock()
		c.rateLimits[rateLimitCategory] = response.Rate
		c.rateMu.Unlock()
	}

	err = CheckResponse(resp)
	if err != nil {
		defer resp.Body.Close()
		// Special case for AcceptedErrors. If an AcceptedError
		// has been encountered, the response's payload will be
		// added to the AcceptedError and returned.
		//
		// Issue #1022
		aerr, ok := err.(*AcceptedError)
		if ok {
			b, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				return response, readErr
			}

			aerr.Raw = b
			err = aerr
		}

		rateLimitError, ok := err.(*RateLimitError)
		if ok && req.Context().Value(SleepUntilPrimaryRateLimitResetWhenRateLimited) != nil {
			if err := sleepUntilResetWithBuffer(req.Context(), rateLimitError.Rate.Reset.Time); err != nil {
				return response, err
			}
			// retry the request once when the rate limit has reset
			return c.BareDo(context.WithValue(req.Context(), SleepUntilPrimaryRateLimitResetWhenRateLimited, nil), req)
		}

		// Update the secondary rate limit if we hit it.
		rerr, ok := err.(*AbuseRateLimitError)
		if ok && rerr.RetryAfter != nil {
			c.rateMu.Lock()
			c.secondaryRateLimitReset = time.Now().Add(*rerr.RetryAfter)
			c.rateMu.Unlock()
		}
	}
	return response, err
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer interface,
// the raw response body will be written to v, without attempting to first
// decode it. If v is nil, and no error happens, the response is returned as is.
// If rate limit is exceeded and reset time is in the future, Do returns
// *RateLimitError immediately without making a network API call.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return resp, err
}