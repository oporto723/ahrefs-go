package ahrefs

import (
	"context"
	"net/http"
)

type ReferringDomainsResponse struct {
	ReferringDomains []ReferringDomain     `json:"refdomains"`
	Stats            ReferringDomainsStats `json:"stats"`
}

type ReferringDomain struct {
	ReferringDomain string `json:"refdomain"`
	DomainRating    int64  `json:"domain_rating"`
	Backlinks       int64  `json:"backlinks"`
}

type ReferringDomainsStats struct {
	ReferringDomains int64 `json:"refdomains"`
	IPs              int64 `json:"ips"`
	ClassC           int64 `json:"class_c"`
}

func (s *serviceImpl) ReferringDomains(ctx context.Context, opts ...Option) (*ReferringDomainsResponse, *http.Response, error) {
	b := s.client.requestBuilder(ctx, opts)
	b.WithColumns("refdomain,domain_rating,backlinks")
	b.WithFrom("refdomains")
	b.WithMode("subdomains")
	b.WithWhere("country=\"us\"")
	b.WithOrderBy("domain_rating:desc")

	payload := &ReferringDomainsResponse{}
	resp, err := s.client.Do(ctx, b, payload)
	if err != nil {
		return nil, resp, err
	}

	return payload, resp, err
}
