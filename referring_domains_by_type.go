package ahrefs

import (
	"context"
	"net/http"
)

type ReferringDomainsByTypeResponse struct {
	ReferringDomainsByType []ReferringDomainByType     `json:"refdomains"`
	Stats                  ReferringDomainsByTypeStats `json:"stats"`
	TLDs                   []TLD                       `json:"tlds"`
}

type ReferringDomainByType struct {
	RefDomain          string  `json:"refdomain"`
	ReferringDomainTop string  `json:"refdomain_top"`
	Backlinks          int64   `json:"backlinks"`
	BacklinksDofollow  int64   `json:"backlinks_dofollow"`
	ReferringPages     int64   `json:"refpages"`
	FirstSeen          string  `json:"first_seen"`
	LastVisited        string  `json:"last_visited"`
	DomainRating       int64   `json:"domain_rating"`
	Traffic            float64 `json:"traffic"`
	LinkedDomains      int64   `json:"linked_domains"`
	RefDomains         int64   `json:"refdomains"`
}

type ReferringDomainsByTypeStats struct {
	RefDomains             int64 `json:"refdomains"`
	IPs                    int64 `json:"ips"`
	ClassC                 int64 `json:"class_c"`
	MaxBacklinks           int64 `json:"max_backlinks"`
	MaxBacklinksDofollow   int64 `json:"max_backlinks_dofollow"`
	MaxRefpages            int64 `json:"max_refpages"`
	TotalBacklinks         int64 `json:"total_backlinks"`
	TotalBacklinksDofollow int64 `json:"total_backlinks_dofollow"`
	All                    int64 `json:"all"`
	Text                   int64 `json:"text"`
	Image                  int64 `json:"image"`
	Nofollow               int64 `json:"nofollow"`
	Ugc                    int64 `json:"ugc"`
	Sponsored              int64 `json:"sponsored"`
	Dofollow               int64 `json:"dofollow"`
	Redirect               int64 `json:"redirect"`
	Canonical              int64 `json:"canonical"`
	Gov                    int64 `json:"gov"`
	Edu                    int64 `json:"edu"`
	RSS                    int64 `json:"rss"`
	Alternate              int64 `json:"alternate"`
}

type TLD struct {
	TLD   string `json:"tld"`
	Count int64  `json:"count"`
}

func (s *serviceImpl) ReferringDomainsByType(ctx context.Context, opts ...Option) (*ReferringDomainsByTypeResponse, *http.Response, error) {
	b := s.client.requestBuilder(ctx, opts)
	b.WithFrom("refdomains_by_type")
	b.WithMode("subdomains")
	b.WithWhere("dofollow=true")
	b.WithOrderBy("domain_rating:desc")

	payload := &ReferringDomainsByTypeResponse{}
	resp, err := s.client.Do(ctx, b, payload)
	if err != nil {
		return nil, resp, err
	}

	return payload, resp, err
}
