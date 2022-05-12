package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/twisp/twisp-sdk-go/pkg/token"
)

var (
	region  = os.Getenv("AWS_REGION")
	account = ""
)

func main() {
	flag.StringVar(&account, "account", "cloud", "which twisp account to use for signing.")
	flag.StringVar(&region, "region", "us-west-2", "the aws region you're authenticating against.")
	flag.Parse()

	url := fmt.Sprintf("https://auth.%s.%s.twisp.com/", region, account)

	g, err := token.NewGenerator(true)
	if err != nil {
		panic(err)
	}

	o := token.GetTokenOptions{
		ClusterID: url,
		Region:    region,
	}

	t, err := g.GetWithOptions(&o)
	if err != nil {
		panic(err)
	}

	j, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
