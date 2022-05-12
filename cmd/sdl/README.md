SDL
-----------

Get the GraphQL SDL from Twisp to use in generating api clients.

## Usage

```
Usage of sdl:
  -account string
    	which twisp account to use for signing. (default "cloud")
  -customer-account string
    	the AWS customer account. If using IAM auth do not set.
  -jwt string
    	an oidc compliant jwt you wish to use. If you use jwt you must specify your aws account.
  -region string
    	the aws region you're authenticating against. (default "us-east-2")
  -schema-out string
    	the location to put the file. If not specified, printed on stdout
```


## Examples

With an IAM credential in the environment:
```
sdl -region=us-west-2 -schema-out schema.graphql
```

With a OIDC compliant JWT:
```
AUTH="<your jwt>"
sdl -region=us-west-2 -customer-account=9999999999 -jwt="$AUTH" -schema-out schema.graphql
```
