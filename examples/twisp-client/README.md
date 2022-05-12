Twisp Client
------------

Uses https://github.com/Khan/genqlient to generate client for Twisp.  This client can then use either IAM token exchange or an OIDC complaint JWT to query Twisp.

Assuming your IAM credentials are in the environment:


## Usage

```
Usage of twisp-client:
  -account string
    	which twisp account to use for signing. (default "cloud")
  -customer-account string
    	the AWS customer account. If using IAM auth do not set.
  -jwt string
    	an oidc compliant jwt you wish to use. If you use jwt you must specify your aws account.
  -region string
    	the aws region you're authenticating against. (default "us-east-2")
```

## Updating Schema

install the `sdl` tool and use update schema from Twisp.
```
go install github.com/twisp/twisp-sdk-go/cmd/sdl
sdl -region=us-west-2 -schema-out=schema.graphql
```

Regenerate client:
```
go generate ./...
```

## Examples

Using IAM credential in environment
```
go install ./...
twisp-client -region=us-west-2
```

With an OIDC compliant JWT:

```
AUTH="<your jwt>"
twisp-client -region=us-west-2 -customer-account=9999999999 -jwt="$AUTH"
```
