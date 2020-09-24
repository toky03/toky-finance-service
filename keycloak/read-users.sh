curl -X GET 'https://id.bubelu.ch/auth/admin/realms/toky-finance/users' \
-H "Accept: application/json" \
-H "Authorization: Bearer $TKN" | jq .
