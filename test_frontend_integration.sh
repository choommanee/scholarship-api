#!/bin/bash

# Test script for Frontend Integration
# Tests Payment and Analytics Services

BASE_URL="http://localhost:8080/api/v1"
ADMIN_EMAIL="admin@university.ac.th"
ADMIN_PASSWORD="admin123"

echo "🧪 Frontend Services Integration Test"
echo "======================================"
echo ""

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counter
PASS=0
FAIL=0

# Test function
test_api() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local headers="$4"
    local expected_status="${5:-200}"

    echo "Testing: $name"

    if [ -z "$headers" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" -H "$headers")
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✅ PASS${NC} - $name (HTTP $http_code)"
        PASS=$((PASS + 1))
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        echo -e "${RED}❌ FAIL${NC} - $name (Expected HTTP $expected_status, got $http_code)"
        FAIL=$((FAIL + 1))
        echo "$body"
    fi
    echo ""
}

# Step 1: Login as admin
echo "🔐 Step 1: Admin Login"
echo "----------------------"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo -e "${RED}❌ Failed to get admin token${NC}"
    echo "$LOGIN_RESPONSE" | jq '.'
    exit 1
fi

echo -e "${GREEN}✅ Admin login successful${NC}"
echo "Token: ${TOKEN:0:20}..."
echo ""

# Step 2: Test Payment Service APIs
echo "📦 Step 2: Testing Payment Service"
echo "-----------------------------------"

test_api \
    "Get Payment Methods" \
    "GET" \
    "/payments/methods" \
    "Authorization: Bearer $TOKEN" \
    200

# Step 3: Test Analytics Service APIs
echo "📊 Step 3: Testing Analytics Service"
echo "-------------------------------------"

test_api \
    "Get Dashboard Summary" \
    "GET" \
    "/analytics/dashboard" \
    "Authorization: Bearer $TOKEN" \
    200

test_api \
    "Get Average Processing Time" \
    "GET" \
    "/analytics/processing-time" \
    "Authorization: Bearer $TOKEN" \
    200

test_api \
    "Get Bottlenecks Analysis" \
    "GET" \
    "/analytics/bottlenecks" \
    "Authorization: Bearer $TOKEN" \
    200

test_api \
    "Get All Statistics" \
    "GET" \
    "/analytics/statistics/all" \
    "Authorization: Bearer $TOKEN" \
    200

# Summary
echo "======================================"
echo "📋 Test Summary"
echo "======================================"
TOTAL=$((PASS + FAIL))
echo -e "Total:  $TOTAL"
echo -e "${GREEN}Passed: $PASS${NC}"
echo -e "${RED}Failed: $FAIL${NC}"
echo "======================================"

if [ $FAIL -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All tests passed!${NC}\n"
    exit 0
else
    echo -e "\n${RED}⚠️  Some tests failed${NC}\n"
    exit 1
fi
