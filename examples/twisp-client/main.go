package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/twisp/twisp-sdk-go/pkg/client"
)

var (
	region          = os.Getenv("AWS_REGION")
	account         = ""
	customerJWT     = ""
	customerAccount = ""
)

type rtWithBilling struct {
	inner http.RoundTripper
}

func (r *rtWithBilling) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Twisp-Include-Billing", "true")
	return r.inner.RoundTrip(req)
}

func main() {
	flag.StringVar(&account, "account", "cloud", "which twisp account to use for signing.")
	flag.StringVar(&region, "region", "us-east-2", "the aws region you're authenticating against.")
	flag.StringVar(&customerAccount, "customer-account", "", "The customer account to target")
	flag.Parse()

	graphqlURL := fmt.Sprintf("https://api.%s.%s.twisp.com/financial/v1/graphql", region, account)

	if customerAccount == "" {
		handle(fmt.Errorf("customer-account is required"))
	}

	twispHTTP := client.NewTwispHttpWithRoundTripper(customerAccount, account, region, &rtWithBilling{inner: http.DefaultTransport})

	// Check a balance
	graphqlClient := client.NewTwispClient(graphqlURL, twispHTTP)
	resp, extBalances, err := CheckAccountBalances(
		context.Background(),
		graphqlClient,
		uuid.MustParse("1fd1dd3e-33fe-4ef5-9d58-676ef8d306b5"),
		uuid.MustParse("822cb59f-ce51-4837-8391-2af3b7a5fc51"),
	)
	handle(err)
	PrintJSON(resp)
	PrintJSON(extBalances)

	txResp, extDeposit, err := PostDeposit(
		context.Background(),
		graphqlClient,
		uuid.Must(uuid.NewRandom()),
		"1fd1dd3e-33fe-4ef5-9d58-676ef8d306b5",
		"1.00",
		"2023-01-01",
	)
	handle(err)
	PrintJSON(txResp)
	PrintJSON(extDeposit)

	UpdateAccountWithOptions(context.Background(), graphqlClient, uuid.New(), AccountUpdateInput{})
}

func PrintJSON(obj any) {
	b, err := json.Marshal(obj)
	handle(err)
	fmt.Printf("%v\n", string(b))
}

func handle(err error) {
	if err != nil {
		log.Fatalf("exiting with error: %s\n", err.Error())
		os.Exit(1)
	}
}

//go:generate go run github.com/twisp/twisp-sdk-go/cmd/generate genqlient.yaml
