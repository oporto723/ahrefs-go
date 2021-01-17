package ahrefs

import (
	"context"
	"net/http"
)

type PositionMetricsResponse struct {
	PositionMetrics PositionMetrics `json:"metrics"`
}

type PositionMetrics struct {
	Positions      int64   `json:"positions"`
	PositionsTop10 int64   `json:"positions_top10"`
	PositionsTop3  int64   `json:"positions_top3"`
	Traffic        float64 `json:"traffic"`
	TrafficTop10   float64 `json:"traffic_top10"`
	TrafficTop3    float64 `json:"traffic_top3"`
	Cost           float64 `json:"cost"`
	CostTop10      float64 `json:"cost_top10"`
	CostTop3       float64 `json:"cost_top3"`
}

func (s *serviceImpl) PositionMetrics(ctx context.Context, opts ...Option) (*PositionMetricsResponse, *http.Response, error) {
	b := s.client.requestBuilder(ctx, opts)
	b.WithFrom("positions_metrics")
	b.WithMode("subdomains")
	b.WithWhere("country=\"us\"")

	payload := &PositionMetricsResponse{}
	resp, err := s.client.Do(ctx, b, payload)
	if err != nil {
		return nil, resp, err
	}

	return payload, resp, err
}
