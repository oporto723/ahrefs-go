package ahrefs

import (
	"context"
	"net/http"
)

type BacklinksOnePerDomainResponse struct {
	Refpages []Refpage `json:"refpages"`
}

type Refpage struct {
	URLFrom          string `json:"url_from"`
	AhrefsRank       int64  `json:"ahrefs_rank"`
	DomainRating     int64  `json:"domain_rating"`
	AhrefsTop        int64  `json:"ahrefs_top"`
	IPFrom           string `json:"ip_from"`
	LinksInternal    int64  `json:"links_internal"`
	LinksExternal    int64  `json:"links_external"`
	PageSize         int64  `json:"page_size"`
	Encoding         string `json:"encoding"`
	Title            string `json:"title"`
	Language         string `json:"language"`
	URLTo            string `json:"url_to"`
	FirstSeen        string `json:"first_seen"`
	LastVisited      string `json:"last_visited"`
	PrevVisited      string `json:"prev_visited"`
	Original         bool   `json:"original"`
	Redirect         int64  `json:"redirect"`
	Anchor           string `json:"anchor"`
	TextPre          string `json:"text_pre"`
	TextPost         string `json:"text_post"`
	HTTPCode         int64  `json:"http_code"`
	URLFromFirstSeen string `json:"url_from_first_seen"`
	FirstOrigin      string `json:"first_origin"`
	LastOrigin       string `json:"last_origin"`
	LinkType         string `json:"link_type"`
	Nofollow         bool   `json:"nofollow"`
	Ugc              bool   `json:"ugc"`
	Sponsored        bool   `json:"sponsored"`
	TotalBacklinks   int64  `json:"total_backlinks"`
}

func (s *serviceImpl) BacklinksOnePerDomain(ctx context.Context, opts ...Option) (*BacklinksOnePerDomainResponse, *http.Response, error) {
	b := s.client.requestBuilder(ctx, opts)
	b.WithFrom("backlinks_one_per_domain")
	b.WithMode("subdomains")

	payload := &BacklinksOnePerDomainResponse{}
	resp, err := s.client.Do(ctx, b, payload)
	if err != nil {
		return nil, resp, err
	}

	return payload, resp, err
}
