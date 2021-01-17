# ahrefs-go

[![GoDoc](https://godoc.org/github.com/oporto723/ahrefs-go?status.svg)](https://godoc.org/github.com/oporto723/ahrefs-go) 
[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Foklog%2Frun%2Fbadge&style=flat-square&label=build)](https://github.com/oporto723/ahrefs-go/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/oporto723/ahrefs-go)](https://goreportcard.com/report/github.com/oporto723/ahrefs-go)

ahrefs-go is a Go client library for accessing the Ahrefs API.

## Usage

```go
client := ahrefs.NewClient(http.DefaultClient, "token")

payload, resp, err := client.Service.ReferringDomainsByType(
    context.TODO(),
    ahrefs.WithTarget("ahrefs.com"),
    ahrefs.WithHaving("domain_rating>10"),
    ahrefs.WithLimit(1))
if err != nil {
    log.Fatal(err)
}s

fmt.Println(resp.StatusCode)
fmt.Println(payload.ReferringDomainsByType)
```