package ahrefs

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://apiv2.ahrefs.com/"
	userAgent      = "go-ahrefs"
)

type Client struct {
	// HTTP client.
	client *http.Client

	// Ahrefs API URL.
	BaseURL *url.URL

	// Client user agent.
	UserAgent string

	// Authencation token.
	token string

	// Service interface.
	Service Service
}

func NewClient(httpClient *http.Client, token string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
		httpClient.Timeout = time.Minute * 2
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
		token:     token,
	}
	c.Service = &serviceImpl{c}

	return c
}

func (c *Client) requestBuilder(ctx context.Context, opts []Option) *requestBuilder {
	return &requestBuilder{
		client: c,
		opts:   opts,
	}
}

// NewRequest prepares an Ahrefs API request. It uses query parameters, as
// defined in https://ahrefs.com/api/documentation. Token and output parameters
// are managed implicitly.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("token")) > 0 {
		params.Set("token", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

func (c *Client) Do(ctx context.Context, builder *requestBuilder, v interface{}) (*http.Response, error) {
	req, err := builder.request()
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req = req.WithContext(ctx)

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

	defer resp.Body.Close()

	if code := resp.StatusCode; code < 200 || code >= 300 {
		return resp, errors.New("unexpected status code")
	}

	// TODO: do not read the whole payload in memory.
	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	if status := resp.Header.Get("X-Status"); status == "error" {
		apiErr := struct {
			Error string `json:"error"`
		}{}
		err = json.Unmarshal(blob, &apiErr)
		if err != nil {
			return resp, err
		}
		return resp, errors.New(apiErr.Error)
	}

	if v == nil {
		return resp, nil
	}

	err = json.Unmarshal(blob, v)
	return resp, err
}

type Service interface {
	ReferringDomains(ctx context.Context, opts ...Option) (*ReferringDomainsResponse, *http.Response, error)
	ReferringDomainsByType(ctx context.Context, opts ...Option) (*ReferringDomainsByTypeResponse, *http.Response, error)
	BacklinksOnePerDomain(ctx context.Context, opts ...Option) (*BacklinksOnePerDomainResponse, *http.Response, error)
	PositionMetrics(ctx context.Context, opts ...Option) (*PositionMetricsResponse, *http.Response, error)
	Pages(ctx context.Context, opts ...Option) (*PagesResponse, *http.Response, error)
}

type serviceImpl struct {
	client *Client
}

var _ Service = &serviceImpl{}
