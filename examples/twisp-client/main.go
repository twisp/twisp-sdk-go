package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/twisp/twisp-sdk-go/pkg/client"
	"github.com/twisp/twisp-sdk-go/pkg/token"
)

var (
	region          = os.Getenv("AWS_REGION")
	account         = ""
	customerJWT     = ""
	customerAccount = ""
)

func main() {
	flag.StringVar(&account, "account", "cloud", "which twisp account to use for signing.")
	flag.StringVar(&region, "region", "us-east-2", "the aws region you're authenticating against.")
	flag.StringVar(&customerJWT, "jwt", "", "an oidc compliant jwt you wish to use. If you use jwt you must specify your aws account.")
	flag.StringVar(&customerAccount, "customer-account", "", "the AWS customer account. If using IAM auth do not set.")
	flag.Parse()

	var isIAM bool
	var graphqlURL string

	if customerJWT != "" && customerAccount == "" {
		handle(fmt.Errorf("customer-account is required"))
	}

	if customerJWT == "" {
		isIAM = true
	}

	var authorization string

	if isIAM {
		graphqlURL = fmt.Sprintf("https://api.%s.%s.twisp.com/graphql", region, account)
		a, err := token.Exchange(account, region)
		handle(err)
		authorization = string(a)
	} else {
		graphqlURL = fmt.Sprintf("https://api.%s.%s.twisp.com/graphql/oidc", region, account)
		authorization = customerJWT
	}

	twispHTTP := client.NewTwispHttp(authorization, customerAccount)

	// Check a balance
	graphqlClient := client.NewTwispClient(graphqlURL, nil, twispHTTP)
	resp, err := checkBalance(
		context.Background(),
		graphqlClient,
		"c9956621-2209-4d0d-bec0-52107fe833fd",
	)
	handle(err)
	PrintJSON(resp)

	//Insert a transaction
	var transactionJSON = `
{
	"account": {
		"account_id": "c9956621-2209-4d0d-bec0-52107fe833fd",
		"status": "OPEN",
		"account_type": "DDA"
	},
	"settlement_account": {
		"account_id": "79109baf-f687-4ccb-b797-c132d64adf36",
		"status": "OPEN",
		"account_type": "GL"
	},
	"transaction_id": "c81505af-f6a2-46cb-8057-7b1f0af3548c",
	"correlation_id": "e7cee9e7-ab07-4030-a470-d3ab05a5fd67",
	"tran_code_id": 3,
	"journal_id": 1,
	"layer_id": 1,
	"effective": "2022-03-20",
	"created": "2022-03-20T16:25:11.000Z",
	"credit": false,
	"amount": 800
}`
	var variables map[string]interface{}
	handle(json.Unmarshal([]byte(transactionJSON), &variables))

	graphqlClient = client.NewTwispClient(graphqlURL, variables, twispHTTP)
	txResp, err := insertTransaction(
		context.Background(),
		graphqlClient,
	)
	handle(err)
	PrintJSON(txResp)

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
