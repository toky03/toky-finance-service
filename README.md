## Configuration Env Variables
- DB_NAME
- DB_PASS
- DB_USER
- DB_TYPE
- DB_HOST
- DB_PORT
- USER_BATCH_PORT
- OPENID_JWKS_URL
#### build browser nginx
`docker build  -t toky03/simpleaccounting-backend .`

## Protocolbuffer
### make protoc-gen-go available
`export GO111MODULE=on`
`export PATH="$PATH:$(go env GOPATH)/bin"`

### Generate the code
`protoc -I grpc_users/ grpc_users/userservice.proto --go_out=plugins=grpc:grpc_users`

