#!/bin/bash

# Test script for Application Details API fixes

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080/api/v1"

echo -e "${YELLOW}Testing Application Details API Fixes${NC}"
echo "========================================"

# Step 1: Login to get token
echo -e "\n${YELLOW}Step 1: Logging in as student...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$TOKEN" ]; then
  echo -e "${RED}❌ Login failed. Response: $LOGIN_RESPONSE${NC}"
  exit 1
fi

echo -e "${GREEN}✓ Login successful${NC}"
echo "Token: ${TOKEN:0:20}..."

# Step 2: Test GetMyApplications endpoint
echo -e "\n${YELLOW}Step 2: Testing GET /applications/my (Problem 1)...${NC}"
MY_APPS_RESPONSE=$(curl -s -X GET "$API_URL/applications/my" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$MY_APPS_RESPONSE" | tail -n 1)
RESPONSE_BODY=$(echo "$MY_APPS_RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "200" ]; then
  echo -e "${GREEN}✓ GetMyApplications works! (200 OK)${NC}"
  echo "Response: $RESPONSE_BODY"
else
  echo -e "${RED}❌ GetMyApplications failed with status $HTTP_CODE${NC}"
  echo "Response: $RESPONSE_BODY"
fi

# Step 3: Create a test application
echo -e "\n${YELLOW}Step 3: Creating a test application...${NC}"
CREATE_APP_RESPONSE=$(curl -s -X POST "$API_URL/applications" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "scholarship_id": 1,
    "family_income": 20000,
    "monthly_expenses": 5000,
    "siblings_count": 2
  }' \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$CREATE_APP_RESPONSE" | tail -n 1)
RESPONSE_BODY=$(echo "$CREATE_APP_RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "201" ]; then
  echo -e "${GREEN}✓ Application created successfully${NC}"
  APPLICATION_ID=$(echo $RESPONSE_BODY | grep -o '"application_id":[0-9]*' | sed 's/"application_id"://')
  if [ -z "$APPLICATION_ID" ]; then
    # Try alternative field name
    APPLICATION_ID=$(echo $RESPONSE_BODY | grep -o '"ApplicationID":[0-9]*' | sed 's/"ApplicationID"://')
  fi
  echo "Application ID: $APPLICATION_ID"
else
  echo -e "${YELLOW}⚠ Application creation returned status $HTTP_CODE (may already exist)${NC}"
  echo "Response: $RESPONSE_BODY"
  # Try to get an existing application ID from my applications
  APPLICATION_ID=$(echo $RESPONSE_BODY | grep -o '"id":"[0-9]*"' | head -1 | sed 's/"id":"//' | sed 's/"//')
fi

# If we still don't have an application ID, try to get from my applications
if [ -z "$APPLICATION_ID" ]; then
  echo -e "${YELLOW}Fetching existing application ID...${NC}"
  APPLICATION_ID=$(echo $RESPONSE_BODY | grep -o '"id":"[0-9]*"' | head -1 | sed 's/"id":"//' | sed 's/"//')
fi

if [ -z "$APPLICATION_ID" ]; then
  echo -e "${RED}❌ Could not get application ID. Using ID=1 for testing.${NC}"
  APPLICATION_ID=1
fi

# Step 4: Test application details endpoint (Problem 2)
echo -e "\n${YELLOW}Step 4: Testing POST /applications/$APPLICATION_ID/personal-info (Problem 2)...${NC}"
PERSONAL_INFO_RESPONSE=$(curl -s -X POST "$API_URL/applications/$APPLICATION_ID/personal-info" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "citizen_id": "1234567890123",
    "date_of_birth": "2000-01-01",
    "gender": "male",
    "nationality": "Thai",
    "phone_number": "0812345678",
    "email": "student@example.com"
  }' \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$PERSONAL_INFO_RESPONSE" | tail -n 1)
RESPONSE_BODY=$(echo "$PERSONAL_INFO_RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "200" ]; then
  echo -e "${GREEN}✓ SavePersonalInfo works! (200 OK)${NC}"
  echo "Response: $RESPONSE_BODY"
elif [ "$HTTP_CODE" = "403" ]; then
  echo -e "${RED}❌ Still getting 403 Forbidden - Problem 2 NOT fixed${NC}"
  echo "Response: $RESPONSE_BODY"
elif [ "$HTTP_CODE" = "404" ]; then
  echo -e "${YELLOW}⚠ Application not found (404) - This is expected if application doesn't exist${NC}"
  echo "Response: $RESPONSE_BODY"
else
  echo -e "${RED}❌ SavePersonalInfo failed with status $HTTP_CODE${NC}"
  echo "Response: $RESPONSE_BODY"
fi

# Step 5: Test another details endpoint
echo -e "\n${YELLOW}Step 5: Testing POST /applications/$APPLICATION_ID/addresses...${NC}"
ADDRESSES_RESPONSE=$(curl -s -X POST "$API_URL/applications/$APPLICATION_ID/addresses" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '[{
    "address_type": "current",
    "house_number": "123",
    "village": "Village 1",
    "subdistrict": "Subdistrict",
    "district": "District",
    "province": "Bangkok",
    "postal_code": "10100"
  }]' \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$ADDRESSES_RESPONSE" | tail -n 1)
RESPONSE_BODY=$(echo "$ADDRESSES_RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "200" ]; then
  echo -e "${GREEN}✓ SaveAddresses works! (200 OK)${NC}"
  echo "Response: $RESPONSE_BODY"
elif [ "$HTTP_CODE" = "403" ]; then
  echo -e "${RED}❌ Still getting 403 Forbidden${NC}"
  echo "Response: $RESPONSE_BODY"
else
  echo -e "${YELLOW}⚠ SaveAddresses returned status $HTTP_CODE${NC}"
  echo "Response: $RESPONSE_BODY"
fi

# Summary
echo -e "\n${YELLOW}========================================${NC}"
echo -e "${YELLOW}Test Summary${NC}"
echo -e "${YELLOW}========================================${NC}"
echo -e "Problem 1 (GetMyApplications 500 error): Check Step 2 result above"
echo -e "Problem 2 (Details endpoints 403 error): Check Steps 4-5 results above"
