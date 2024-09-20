# Requirements for local Dev setup
- Build keycloak from keycloak_example

## Configuration Env Variables
- DB_NAME
- DB_PASS
- DB_USER
- DB_TYPE
- DB_HOST
- DB_PORT
- USER_BATCH_PORT
- OPENID_JWKS_URL
- ID_PROVIDER_CLIENT_SECRET can be found in keycloack under Clients > Client details > Credentials
- OPENID_JWKS_EXTERNAL_URL optional if not set it will be the same as OPENID_JWKS_URL
#### build
`docker build -t toky03/simpleaccounting-backend .`

## Protocolbuffer
### make protoc-gen-go available
`export GO111MODULE=on`
`export PATH="$PATH:$(go env GOPATH)/bin"`

### Generate the code from withing root directory
`protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc_users/userservice.proto`

## Saved testuser
Username: toky
Password: pwd

