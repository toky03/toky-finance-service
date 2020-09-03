## Configuration Env Variables
- USER_BATCH_PORT
- OPENID_JWKS_URL

## Protocolbuffer
### make protoc-gen-go available
`export GO111MODULE=on`
`export PATH="$PATH:$(go env GOPATH)/bin"`

### Generate the code
`protoc -I grpc_users/ grpc_users/userservice.proto --go_out=plugins=grpc:grpc_users`

