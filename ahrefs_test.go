package ahrefs_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/oporto723/ahrefs-go"
)

// setup is our helper to create a new ahrefs client with a fakeserver.
func setup(t *testing.T, fakeServer http.Handler) *ahrefs.Client {
	t.Helper()

	srv := httptest.NewServer(fakeServer)

	client := ahrefs.NewClient(nil, "12345")
	url, _ := url.Parse(srv.URL + "/")
	client.BaseURL = url

	t.Cleanup(func() {
		srv.Close()
	})

	return client
}

func TestRealAPI(t *testing.T) {
	t.Skip("for demonstration purposes")

	t.Parallel()
	c := qt.New(t)

	const demoToken = "2bc0fc25a447ff7ff1fc8998e4cd2eeb8cb5fc7e"
	client := ahrefs.NewClient(nil, demoToken)

	ctx := context.Background()
	payload, resp, err := client.Service.ReferringDomainsByType(ctx, ahrefs.WithTarget("ahrefs.com"), ahrefs.WithLimit(1))

	c.Assert(err, qt.IsNil)
	c.Assert(resp.StatusCode, qt.Equals, http.StatusOK)
	c.Assert(payload, qt.Not(qt.IsNil))
	c.Assert(len(payload.ReferringDomainsByType), qt.Equals, 1)
}

func TestQueryParameters(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	fakeServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Results-Count", "0")
		w.Header().Set("X-Status", "error")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"error":"invalid token"}`))

		queryValues := r.URL.Query()
		c.Assert(queryValues, qt.DeepEquals, url.Values{
			"where":  {"country=\"us\""},
			"from":   {"positions_metrics"},
			"limit":  {"1"},
			"mode":   {"subdomains"},
			"output": {"json"},
			"target": {"ahrefs.com"},
			"token":  {"12345"},
		})
	})

	client := setup(t, fakeServer)

	ctx := context.Background()
	payload, resp, err := client.Service.PositionMetrics(ctx, ahrefs.WithTarget("ahrefs.com"), ahrefs.WithLimit(1))

	c.Assert(err, qt.ErrorMatches, "invalid token")
	c.Assert(resp.StatusCode, qt.Equals, http.StatusOK)
	c.Assert(payload, qt.IsNil)
}

func TestAPIError(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	fakeServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Results-Count", "0")
		w.Header().Set("X-Status", "error")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"error":"order_by: column 'foobar_rating' not found"}`))
	})

	client := setup(t, fakeServer)

	ctx := context.Background()
	payload, resp, err := client.Service.PositionMetrics(ctx, ahrefs.WithTarget("ahrefs.com"), ahrefs.WithLimit(1))

	c.Assert(err, qt.ErrorMatches, "order_by: column 'foobar_rating' not found")
	c.Assert(resp.StatusCode, qt.Equals, http.StatusOK)
	c.Assert(payload, qt.IsNil)
}

func TestReferringDomainsByType(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	fakeserver := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Results-Count", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"refdomains": [
			  {
				"refdomain": "linkedin.com",
				"refdomain_top": "linkedin.com",
				"backlinks": 3,
				"backlinks_dofollow": 0,
				"refpages": 3,
				"first_seen": "2020-01-14T05:16:39Z",
				"last_visited": "2020-12-09T09:54:58Z",
				"domain_rating": 98,
				"traffic": 163250217.345815,
				"linked_domains": 903,
				"refdomains": 8724253
			  }
			],
			"stats": {
			  "refdomains": 48974,
			  "ips": 37383,
			  "class_c": 18533,
			  "max_backlinks": 656327,
			  "max_backlinks_dofollow": 656327,
			  "max_refpages": 551209,
			  "total_backlinks": 3947793,
			  "total_backlinks_dofollow": 3193038,
			  "all": 48974,
			  "text": 51145,
			  "image": 2226,
			  "nofollow": 12325,
			  "ugc": 617,
			  "sponsored": 66,
			  "dofollow": 38754,
			  "redirect": 375,
			  "canonical": 35,
			  "gov": 13,
			  "edu": 266,
			  "rss": 45,
			  "alternate": 6
			},
			"tlds": [
			  {
				"tld": "com",
				"count": 29534
			  },
			  {
				"tld": "de",
				"count": 1901
			  },
			  {
				"tld": "net",
				"count": 1810
			  }
			]
		  }`))
	})

	client := setup(t, fakeserver)

	ctx := context.Background()
	payload, resp, err := client.Service.ReferringDomainsByType(ctx, ahrefs.WithTarget("ahrefs.com"), ahrefs.WithLimit(1))

	c.Assert(err, qt.IsNil)
	c.Assert(resp.StatusCode, qt.Equals, http.StatusOK)
	c.Assert(payload, qt.DeepEquals, &ahrefs.ReferringDomainsByTypeResponse{
		ReferringDomainsByType: []ahrefs.ReferringDomainByType{
			{
				RefDomain:          "linkedin.com",
				ReferringDomainTop: "linkedin.com",
				Backlinks:          3,
				ReferringPages:     3,
				FirstSeen:          "2020-01-14T05:16:39Z",
				LastVisited:        "2020-12-09T09:54:58Z",
				DomainRating:       98,
				Traffic:            1.63250217345815e+08,
				LinkedDomains:      903,
				RefDomains:         8.724253e+06,
			},
		},
		Stats: ahrefs.ReferringDomainsByTypeStats{
			RefDomains:             48974,
			IPs:                    37383,
			ClassC:                 18533,
			MaxBacklinks:           656327,
			MaxBacklinksDofollow:   656327,
			MaxRefpages:            551209,
			TotalBacklinks:         3947793,
			TotalBacklinksDofollow: 3193038,
			All:                    48974,
			Text:                   51145,
			Image:                  2226,
			Nofollow:               12325,
			Ugc:                    617,
			Sponsored:              66,
			Dofollow:               38754,
			Redirect:               375,
			Canonical:              35,
			Gov:                    13,
			Edu:                    266,
			RSS:                    45,
			Alternate:              6,
		},
		TLDs: []ahrefs.TLD{
			{
				TLD:   "com",
				Count: 29534,
			},
			{
				TLD:   "de",
				Count: 1901,
			},
			{
				TLD:   "net",
				Count: 1810,
			},
		},
	})
}

func TestBacklinksOnePerDomain(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	fakeserver := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Results-Count", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"refpages": [
			  {
				"url_from": "https://worldwidetopsiteslist.blogspot.com/2020/10/theworldsmostvisitedwebpages203.html",
				"ahrefs_rank": 5,
				"domain_rating": 2,
				"ahrefs_top": 46372866,
				"ip_from": "216.58.198.193",
				"links_internal": 21,
				"links_external": 7463,
				"page_size": 84102,
				"encoding": "utf8",
				"title": "World Wide top Sites List : the_worlds_most_visited_web_pages_203",
				"language": "en",
				"url_to": "https://www.directcbdonline.com/",
				"first_seen": "2020-10-28T12:23:00Z",
				"last_visited": "2020-12-07T21:53:23Z",
				"prev_visited": "2020-12-01T05:34:47Z",
				"original": true,
				"redirect": 0,
				"alt": "",
				"anchor": "https://directcbdonline.com",
				"text_pre": "572086 https://betongmaidanang.com 572087 https://planetel.it 572088",
				"text_post": "572089 https://yeahfieldtrip.com 572090 https://leonardoborges.com 572091",
				"http_code": 200,
				"url_from_first_seen": "2020-10-25T22:24:38Z",
				"first_origin": "fresh",
				"last_origin": "recrawl",
				"link_type": "href",
				"nofollow": false,
				"ugc": false,
				"sponsored": false,
				"total_backlinks": 2
			  }
			]
		  }`))
	})

	client := setup(t, fakeserver)

	ctx := context.Background()
	payload, resp, err := client.Service.BacklinksOnePerDomain(ctx, ahrefs.WithTarget("ahrefs.com"), ahrefs.WithLimit(1))

	c.Assert(err, qt.IsNil)
	c.Assert(resp.StatusCode, qt.Equals, http.StatusOK)
	c.Assert(payload, qt.DeepEquals, &ahrefs.BacklinksOnePerDomainResponse{
		Refpages: []ahrefs.Refpage{
			{
				URLFrom:          "https://worldwidetopsiteslist.blogspot.com/2020/10/theworldsmostvisitedwebpages203.html",
				AhrefsRank:       5,
				DomainRating:     2,
				AhrefsTop:        46372866,
				IPFrom:           "216.58.198.193",
				LinksInternal:    21,
				LinksExternal:    7463,
				PageSize:         84102,
				Encoding:         "utf8",
				Title:            "World Wide top Sites List : the_worlds_most_visited_web_pages_203",
				Language:         "en",
				URLTo:            "https://www.directcbdonline.com/",
				FirstSeen:        "2020-10-28T12:23:00Z",
				LastVisited:      "2020-12-07T21:53:23Z",
				PrevVisited:      "2020-12-01T05:34:47Z",
				Original:         true,
				Redirect:         0,
				Anchor:           "https://directcbdonline.com",
				TextPre:          "572086 https://betongmaidanang.com 572087 https://planetel.it 572088",
				TextPost:         "572089 https://yeahfieldtrip.com 572090 https://leonardoborges.com 572091",
				HTTPCode:         200,
				URLFromFirstSeen: "2020-10-25T22:24:38Z",
				FirstOrigin:      "fresh",
				LastOrigin:       "recrawl",
				LinkType:         "href",
				Nofollow:         false,
				Ugc:              false,
				Sponsored:        false,
				TotalBacklinks:   2,
			},
		},
	})
}

func TestPositionMetrics(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	fakeServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Results-Count", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
		  "metrics": {
			"positions": 98943,
			"positions_top10": 23058,
			"positions_top3": 7131,
			"traffic": 285788.195958,
			"traffic_top10": 273514.303273,
			"traffic_top3": 199336.595846,
			"cost": 121305249.158165,
			"cost_top10": 115373162.606672,
			"cost_top3": 91302494.31712399
		  }
		}`))
	})

	client := setup(t, fakeServer)

	ctx := context.Background()
	payload, resp, err := client.Service.PositionMetrics(ctx, ahrefs.WithTarget("ahrefs.com"), ahrefs.WithLimit(1))

	c.Assert(err, qt.IsNil)
	c.Assert(resp.StatusCode, qt.Equals, http.StatusOK)
	c.Assert(payload, qt.DeepEquals, &ahrefs.PositionMetricsResponse{
		PositionMetrics: ahrefs.PositionMetrics{
			Positions:      98943,
			PositionsTop10: 23058,
			PositionsTop3:  7131,
			Traffic:        285788.195958,
			TrafficTop10:   273514.303273,
			TrafficTop3:    199336.595846,
			Cost:           121305249.158165,
			CostTop10:      115373162.606672,
			CostTop3:       91302494.31712399,
		},
	})
}

func TestPages(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	fakeserver := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Results-Count", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"pages": [{
				"url": "https://petlifetoday.com/",
				"ahrefs_rank": 83,
				"first_seen": "2018-12-20T17:01:04Z",
				"last_visited": "2020-12-15T03:44:34Z",
				"http_code": 200,
				"size": 29239,
				"links_internal": 173,
				"links_external": 3,
				"encoding": "utf8",
				"title": "Pet Care Advice, Tips and Product Reviews - Pet Life Today",
				"redirect_url": "",
				"content_encoding": "br"
			}],
			"stats": {
				"pages": 2991
			}
		}`))
	})

	client := setup(t, fakeserver)

	ctx := context.Background()
	payload, resp, err := client.Service.Pages(ctx, ahrefs.WithTarget("ahrefs.com"), ahrefs.WithLimit(1))

	c.Assert(err, qt.IsNil)
	c.Assert(resp.StatusCode, qt.Equals, http.StatusOK)
	c.Assert(payload, qt.DeepEquals, &ahrefs.PagesResponse{
		Pages: []ahrefs.Page{
			{
				URL:             "https://petlifetoday.com/",
				AhrefsRank:      83,
				FirstSeen:       "2018-12-20T17:01:04Z",
				LastVisited:     "2020-12-15T03:44:34Z",
				HTTPCode:        200,
				Size:            29239,
				LinksInternal:   173,
				LinksExternal:   3,
				Encoding:        "utf8",
				Title:           "Pet Care Advice, Tips and Product Reviews - Pet Life Today",
				RedirectURL:     "",
				ContentEncoding: "br",
			},
		},
		Stats: ahrefs.PagesStats{
			Pages: 2991,
		},
	})
}
