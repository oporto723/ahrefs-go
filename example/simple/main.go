package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/oporto723/ahrefs-go"
)

func main() {
	const demoToken = "2bc0fc25a447ff7ff1fc8998e4cd2eeb8cb5fc7e"
	client := ahrefs.NewClient(http.DefaultClient, demoToken)

	payload, resp, err := client.Service.ReferringDomainsByType(
		context.TODO(),
		ahrefs.WithTarget("ahrefs.com"),
		ahrefs.WithHaving("domain_rating>10"),
		ahrefs.WithLimit(1))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(payload.ReferringDomainsByType)
}
