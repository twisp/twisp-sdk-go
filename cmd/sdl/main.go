package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/twisp/twisp-sdk-go/pkg/client"
)

var (
	region          = os.Getenv("AWS_REGION")
	account         = ""
	customerJWT     = ""
	customerAccount = ""
	schemaOut       = ""
	c               *http.Client
)

func main() {

	flag.StringVar(&account, "account", "cloud", "which twisp account to use for signing.")
	flag.StringVar(&region, "region", "us-east-2", "the aws region you're authenticating against.")
	flag.StringVar(&customerAccount, "customer-account", "", "the customer account id.")

	flag.StringVar(&schemaOut, "schema-out", "", "the location to put the file. If not specified, printed on stdout")
	flag.Parse()

	var graphqlURL string

	if customerAccount == "" {
		handle(fmt.Errorf("customer-account is required"))
	}

	graphqlURL = fmt.Sprintf("https://api.%s.%s.twisp.com/financial/v1/graphql", region, account)
	c = client.NewTwispHttp(customerAccount, account, region)

	q := Query{
		Query:         `{ _service { sdl } }`,
		Variables:     map[string]any{},
		OperationName: "",
	}

	query, err := json.Marshal(q)
	handle(err)

	req, err := http.NewRequest(http.MethodPost, graphqlURL, bytes.NewBuffer(query))
	handle(err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := c.Do(req)
	handle(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handle(err)

	if resp.StatusCode/100 != 2 {
		handle(fmt.Errorf("introspection got response %s %d %s", req.URL.String(), resp.StatusCode, string(body)))
	}

	var bodyJSON map[string]any
	handle(json.Unmarshal(body, &bodyJSON))

	var sdl string
	data := bodyJSON["data"].(map[string]any)
	_service := data["_service"].(map[string]any)
	sdl = _service["sdl"].(string)

	if schemaOut == "" {
		fmt.Printf("%s\n", sdl)
	} else {
		handle(ioutil.WriteFile(schemaOut, []byte(sdl), 0644))
	}
}

func handle(err error) {
	if err != nil {
		log.Fatalf("exiting with error: %s\n", err.Error())
		os.Exit(1)
	}
}

type Query struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName"`
	Variables     map[string]any `json:"variables"`
}
