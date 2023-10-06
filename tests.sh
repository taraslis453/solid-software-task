#!/bin/bash

BASE_URL="http://localhost:8080"

# Register User
echo "Testing User Registration..."
curl -X POST $BASE_URL/users/register/ -H "Content-Type: application/json" -d '{
    "name": "John",
    "surname": "Doe",
    "email": "john.doe@example.com",
    "password": "password123",
    "phone": "+1234567890"
}'

echo -e "\n\n"

# Login User
echo "Testing User Login..."
LOGIN_RESPONSE=$(curl -s -X GET "$BASE_URL/users/login?email=john.doe@example.com&password=password123")
echo $LOGIN_RESPONSE

ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r .accessToken)
USER_ID=$(echo $LOGIN_RESPONSE | jq -r .userId)

echo -e "\n\n"

# Refresh Token
echo "Testing Refresh Token..."
curl -X POST $BASE_URL/users/refresh-token -H "Authorization: Bearer $ACCESS_TOKEN"

echo -e "\n\n"

# Get User Information
echo "Fetching User Information..."
curl -X GET $BASE_URL/users/$USER_ID -H "Authorization: Bearer $ACCESS_TOKEN"

echo -e "\n\n"

# Update User Information
echo "Testing User Update..."
curl -X PUT $BASE_URL/users/ -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{
    "user": {
        "id": "'$USER_ID'",
        "name": "John Updated",
        "surname": "Doe Updated"
    }
}'

echo -e "\n\n"

# Invalid Token
echo "Testing with Invalid Token..."
curl -X GET $BASE_URL/users/$USER_ID -H "Authorization: Bearer InvalidTokenHere"
echo -e "\n\n"

# Invalid Password
echo "Testing with Invalid Password..."
curl -X GET "$BASE_URL/users/login?email=john.doe@example.com&password=wrongpassword"
echo -e "\n\n"

# User Not Found
echo "Testing with Non-existent User..."
NON_EXISTENT_USER_ID="3f096523-8969-4b8e-b48d-b99c1600de9d"
curl -X GET $BASE_URL/users/$NON_EXISTENT_USER_ID -H "Authorization: Bearer $ACCESS_TOKEN"

echo -e "\n\n"
echo "Tests Completed!"

