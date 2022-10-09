Run Keycloak in docker 
` 9225  docker run -p 8080:8080 -e KEYCLOAK_ADMIN=admin --mount type=bind,source=$(pwd),target=/opt/keycloak/themes/accounting-theme -e KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak start-dev`
`