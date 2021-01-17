package ahrefs

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type requestBuilder struct {
	client *Client

	columns string
	from    string
	target  string
	mode    string
	where   string
	having  string
	orderBy string
	limit   string

	// User-provided options.
	opts []Option
}

func (r *requestBuilder) WithColumns(columns string) *requestBuilder {
	r.columns = columns
	return r
}

func (r *requestBuilder) WithFrom(from string) *requestBuilder {
	r.from = from
	return r
}

func (r *requestBuilder) WithTarget(target string) *requestBuilder {
	r.target = target
	return r
}

func (r *requestBuilder) WithMode(mode string) *requestBuilder {
	r.mode = mode
	return r
}

func (r *requestBuilder) WithWhere(where string) *requestBuilder {
	r.where = where
	return r
}

func (r *requestBuilder) WithHaving(having string) *requestBuilder {
	r.having = having
	return r
}

func (r *requestBuilder) WithOrderBy(orderBy string) *requestBuilder {
	r.orderBy = orderBy
	return r
}

func (r *requestBuilder) WithLimit(limit string) *requestBuilder {
	r.limit = limit
	return r
}

func (r *requestBuilder) request() (*http.Request, error) {
	if !strings.HasSuffix(r.client.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", r.client.BaseURL)
	}
	const urlStr = "/"
	u, err := r.client.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	for _, fn := range r.opts {
		fn(r)
	}

	q.Add("output", "json")
	if r.client.token != "" {
		q.Add("token", r.client.token)
	}

	if r.columns != "" {
		q.Add("select", r.columns)
	}
	if r.from != "" {
		q.Add("from", r.from)
	}
	if r.target != "" {
		q.Add("target", r.target)
	}
	if r.mode != "" {
		q.Add("mode", r.mode)
	}
	if r.where != "" {
		q.Add("where", r.where)
	}
	if r.having != "" {
		q.Add("having", r.having)
	}
	if r.orderBy != "" {
		q.Add("orderBy", r.orderBy)
	}
	if r.limit != "" {
		q.Add("limit", r.limit)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if r.client.UserAgent != "" {
		req.Header.Set("User-Agent", r.client.UserAgent)
	}

	return req, nil
}

// Option is a function that changes the request.
type Option func(*requestBuilder)

func WithTarget(target string) Option {
	return func(rb *requestBuilder) {
		rb.WithTarget(target)
	}
}

func WithLimit(limit int64) Option {
	return func(rb *requestBuilder) {
		rb.WithLimit(strconv.FormatInt(limit, 10))
	}
}

func WithHaving(having string) Option {
	return func(rb *requestBuilder) {
		rb.WithHaving(having)
	}
}

func WithWhere(where string) Option {
	return func(rb *requestBuilder) {
		rb.WithWhere(where)
	}
}
