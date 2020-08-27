curl -X GET 'http://localhost:8080/auth/admin/realms/toky-finance/users' \
-H "Accept: application/json" \
-H "Authorization: Bearer $TKN" | jq .
