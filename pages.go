package ahrefs

import (
	"context"
	"net/http"
)

type PagesResponse struct {
	Pages []Page     `json:"pages"`
	Stats PagesStats `json:"stats"`
}

type Page struct {
	URL             string `json:"url"`
	AhrefsRank      int64  `json:"ahrefs_rank"`
	FirstSeen       string `json:"first_seen"`
	LastVisited     string `json:"last_visited"`
	HTTPCode        int64  `json:"http_code"`
	Size            int64  `json:"size"`
	LinksInternal   int64  `json:"links_internal"`
	LinksExternal   int64  `json:"links_external"`
	Encoding        string `json:"encoding"`
	Title           string `json:"title"`
	RedirectURL     string `json:"redirect_url"`
	ContentEncoding string `json:"content_encoding"`
}

type PagesStats struct {
	Pages int64 `json:"pages"`
}

func (s *serviceImpl) Pages(ctx context.Context, opts ...Option) (*PagesResponse, *http.Response, error) {
	b := s.client.requestBuilder(ctx, opts)
	b.WithFrom("pages")
	b.WithMode("subdomains")

	payload := &PagesResponse{}
	resp, err := s.client.Do(ctx, b, payload)
	if err != nil {
		return nil, resp, err
	}

	return payload, resp, err
}
