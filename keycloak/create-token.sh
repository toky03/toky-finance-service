export TKN=$(curl -X POST 'http://localhost:8080/auth/realms/toky-finance/protocol/openid-connect/token' \
 -H "Content-Type: application/x-www-form-urlencoded" \
 -d "client_secret=46dd09e3-19f6-457d-8ec3-0ea0264cefde" \
 -d "username=finance-service" \
 -d 'password=admin' \
 -d 'grant_type=password' \
  -d 'scope=openid' \
 -d 'client_id=toky-finance-service' | jq -r '.access_token')

#   -d "client_secret=46dd09e3-19f6-457d-8ec3-0ea0264cefde" \
