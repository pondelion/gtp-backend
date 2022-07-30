set -u
curl "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=${GCP_IP_API_KEY}" -H 'Content-Type: application/json' --data-binary "{\"email\":\"${GCP_IP_AUTH_TEST_USER_EMAIL}\",\"password\":\"${GCP_IP_AUTH_TEST_USER_PASSWORD}\",\"returnSecureToken\":true}" | jq