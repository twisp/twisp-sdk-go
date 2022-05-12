Twisp Client
------------

Uses https://github.com/Khan/genqlient to generate client for Twisp.  This client can then use either IAM token exchange or an OIDC complaint JWT to query Twisp.

Assuming your IAM credentials are in the environment:

update schema:
```
go install github.com/twisp/twisp-sdk-go/cmd/sdl
./sdl -region="<your-aws-region>" -schema-out="$(PWD)/schema.graphql"
```

Regenerate client:
```
go generate ./...
```

Run:
```
go install ./...
./twisp-client -region="<your-aws-region>"
```
