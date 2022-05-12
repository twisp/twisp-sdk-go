package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Khan/genqlient/graphql"
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

	graphqlClient := graphql.NewClient(graphqlURL, client.NewTwispClient(authorization, customerAccount))

	resp, err := checkBalance(context.Background(), graphqlClient, "c9956621-2209-4d0d-bec0-52107fe833fd")

	handle(err)
	b, err := json.Marshal(resp)
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
