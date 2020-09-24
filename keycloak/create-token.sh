export TKN=$(curl -X POST 'https://id.bubelu.ch/auth/realms/toky-finance/protocol/openid-connect/token' \
 -H "Content-Type: application/x-www-form-urlencoded" \
 -d "client_secret=f5950669-4310-476b-a608-031cb2dbd0bc" \
 -d "username=finance-service" \
 -d 'password=admin' \
 -d 'grant_type=password' \
  -d 'scope=openid' \
 -d 'client_id=toky-finance-service' | jq -r '.access_token')

#   -d "client_secret=46dd09e3-19f6-457d-8ec3-0ea0264cefde" \
