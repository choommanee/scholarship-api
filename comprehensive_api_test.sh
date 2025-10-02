#!/bin/bash

# Comprehensive API Test Script
# ทดสอบ API ทั้งหมดแบบละเอียด พร้อมบันทึกทุกฟิลด์

BASE_URL="http://localhost:8080/api/v1"
OUTPUT_FILE="comprehensive_test_results.json"
REPORT_FILE="COMPREHENSIVE_TEST_REPORT.md"

# Test credentials
STUDENT_EMAIL="test@example.com"
STUDENT_PASSWORD="password123"
ADMIN_EMAIL="admin@university.ac.th"
ADMIN_PASSWORD="admin123"

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Initialize test results
echo "{" > "$OUTPUT_FILE"
echo "  \"test_date\": \"$(date '+%Y-%m-%d %H:%M:%S')\"," >> "$OUTPUT_FILE"
echo "  \"base_url\": \"$BASE_URL\"," >> "$OUTPUT_FILE"
echo "  \"results\": [" >> "$OUTPUT_FILE"

FIRST_TEST=true

# Function to add test result
add_result() {
    local category="$1"
    local name="$2"
    local method="$3"
    local endpoint="$4"
    local status_code="$5"
    local response="$6"
    local expected="$7"
    local result="$8"

    if [ "$FIRST_TEST" = false ]; then
        echo "," >> "$OUTPUT_FILE"
    fi
    FIRST_TEST=false

    cat >> "$OUTPUT_FILE" <<EOF
    {
      "category": "$category",
      "test_name": "$name",
      "method": "$method",
      "endpoint": "$endpoint",
      "expected_status": $expected,
      "actual_status": $status_code,
      "result": "$result",
      "response": $response
    }
EOF
}

# Counter
TOTAL=0
PASS=0
FAIL=0

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║       COMPREHENSIVE API TEST - SCHOLARSHIP SYSTEM              ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

#==============================================================================
# 1. HEALTH CHECK
#==============================================================================
echo -e "${BLUE}┌─────────────────────────────────────────────────────────────┐${NC}"
echo -e "${BLUE}│ 1. HEALTH CHECK & SYSTEM STATUS                            │${NC}"
echo -e "${BLUE}└─────────────────────────────────────────────────────────────┘${NC}"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/../health")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

echo "Test: Health Check"
echo "Endpoint: GET /health"
if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    add_result "System" "Health Check" "GET" "/health" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "System" "Health Check" "GET" "/health" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null || echo "$body"
echo ""

#==============================================================================
# 2. AUTHENTICATION
#==============================================================================
echo -e "${BLUE}┌─────────────────────────────────────────────────────────────┐${NC}"
echo -e "${BLUE}│ 2. AUTHENTICATION & AUTHORIZATION                          │${NC}"
echo -e "${BLUE}└─────────────────────────────────────────────────────────────┘${NC}"

# 2.1 Student Login
echo "Test: Student Login"
echo "Endpoint: POST /auth/login"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$STUDENT_EMAIL\",\"password\":\"$STUDENT_PASSWORD\"}")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    STUDENT_TOKEN=$(echo "$body" | jq -r '.token')
    add_result "Authentication" "Student Login" "POST" "/auth/login" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
    echo "Token: ${STUDENT_TOKEN:0:30}..."
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "Authentication" "Student Login" "POST" "/auth/login" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

# 2.2 Admin Login
echo "Test: Admin Login"
echo "Endpoint: POST /auth/login"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    ADMIN_TOKEN=$(echo "$body" | jq -r '.token')
    add_result "Authentication" "Admin Login" "POST" "/auth/login" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
    echo "Token: ${ADMIN_TOKEN:0:30}..."
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "Authentication" "Admin Login" "POST" "/auth/login" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

#==============================================================================
# 3. USER PROFILE
#==============================================================================
echo -e "${BLUE}┌─────────────────────────────────────────────────────────────┐${NC}"
echo -e "${BLUE}│ 3. USER PROFILE MANAGEMENT                                  │${NC}"
echo -e "${BLUE}└─────────────────────────────────────────────────────────────┘${NC}"

echo "Test: Get User Profile (Student)"
echo "Endpoint: GET /user/profile"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/user/profile" \
    -H "Authorization: Bearer $STUDENT_TOKEN")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    add_result "User" "Get User Profile" "GET" "/user/profile" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "User" "Get User Profile" "GET" "/user/profile" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

#==============================================================================
# 4. SCHOLARSHIPS
#==============================================================================
echo -e "${BLUE}┌─────────────────────────────────────────────────────────────┐${NC}"
echo -e "${BLUE}│ 4. SCHOLARSHIP MANAGEMENT                                   │${NC}"
echo -e "${BLUE}└─────────────────────────────────────────────────────────────┘${NC}"

# 4.1 Get All Scholarships
echo "Test: Get All Scholarships (Public)"
echo "Endpoint: GET /scholarships"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/scholarships")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    count=$(echo "$body" | jq '.total // 0')
    add_result "Scholarship" "Get All Scholarships" "GET" "/scholarships" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
    echo "Total scholarships: $count"
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "Scholarship" "Get All Scholarships" "GET" "/scholarships" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

# 4.2 Get Available Scholarships
echo "Test: Get Available Scholarships"
echo "Endpoint: GET /scholarships/available"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/scholarships/available")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    add_result "Scholarship" "Get Available Scholarships" "GET" "/scholarships/available" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "Scholarship" "Get Available Scholarships" "GET" "/scholarships/available" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

# 4.3 Get Scholarship by ID
SCHOLARSHIP_ID=$(curl -s "$BASE_URL/scholarships" | jq -r '.scholarships[0].scholarship_id // empty')
if [ -n "$SCHOLARSHIP_ID" ]; then
    echo "Test: Get Scholarship by ID"
    echo "Endpoint: GET /scholarships/$SCHOLARSHIP_ID"
    response=$(curl -s -w "\n%{http_code}" "$BASE_URL/scholarships/$SCHOLARSHIP_ID")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    TOTAL=$((TOTAL + 1))

    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
        PASS=$((PASS + 1))
        add_result "Scholarship" "Get Scholarship by ID" "GET" "/scholarships/{id}" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
    else
        echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
        FAIL=$((FAIL + 1))
        add_result "Scholarship" "Get Scholarship by ID" "GET" "/scholarships/{id}" "$http_code" "\"$body\"" "200" "FAIL"
    fi
    echo "$body" | jq '.' 2>/dev/null
    echo ""
fi

#==============================================================================
# 5. NEWS
#==============================================================================
echo -e "${BLUE}┌─────────────────────────────────────────────────────────────┐${NC}"
echo -e "${BLUE}│ 5. NEWS & ANNOUNCEMENTS                                     │${NC}"
echo -e "${BLUE}└─────────────────────────────────────────────────────────────┘${NC}"

# 5.1 Get All News
echo "Test: Get All News (Public)"
echo "Endpoint: GET /news"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/news")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    add_result "News" "Get All News" "GET" "/news" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "News" "Get All News" "GET" "/news" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

#==============================================================================
# 6. PAYMENT SYSTEM (Admin)
#==============================================================================
echo -e "${BLUE}┌─────────────────────────────────────────────────────────────┐${NC}"
echo -e "${BLUE}│ 6. PAYMENT SYSTEM (Admin Access)                           │${NC}"
echo -e "${BLUE}└─────────────────────────────────────────────────────────────┘${NC}"

# 6.1 Get Payment Methods
echo "Test: Get Payment Methods"
echo "Endpoint: GET /payments/methods"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/payments/methods" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    add_result "Payment" "Get Payment Methods" "GET" "/payments/methods" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "Payment" "Get Payment Methods" "GET" "/payments/methods" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

#==============================================================================
# 7. ANALYTICS SYSTEM (Admin)
#==============================================================================
echo -e "${BLUE}┌─────────────────────────────────────────────────────────────┐${NC}"
echo -e "${BLUE}│ 7. ANALYTICS & REPORTING (Admin Access)                    │${NC}"
echo -e "${BLUE}└─────────────────────────────────────────────────────────────┘${NC}"

# 7.1 Dashboard Summary
echo "Test: Get Dashboard Summary"
echo "Endpoint: GET /analytics/dashboard"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/analytics/dashboard" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    add_result "Analytics" "Dashboard Summary" "GET" "/analytics/dashboard" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "Analytics" "Dashboard Summary" "GET" "/analytics/dashboard" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

# 7.2 Processing Time
echo "Test: Get Average Processing Time"
echo "Endpoint: GET /analytics/processing-time"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/analytics/processing-time" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    add_result "Analytics" "Processing Time" "GET" "/analytics/processing-time" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "Analytics" "Processing Time" "GET" "/analytics/processing-time" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

# 7.3 Bottlenecks
echo "Test: Get Bottlenecks Analysis"
echo "Endpoint: GET /analytics/bottlenecks"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/analytics/bottlenecks" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    add_result "Analytics" "Bottlenecks Analysis" "GET" "/analytics/bottlenecks" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "Analytics" "Bottlenecks Analysis" "GET" "/analytics/bottlenecks" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

# 7.4 All Statistics
echo "Test: Get All Statistics"
echo "Endpoint: GET /analytics/statistics/all"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/analytics/statistics/all" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
TOTAL=$((TOTAL + 1))

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $http_code)"
    PASS=$((PASS + 1))
    add_result "Analytics" "All Statistics" "GET" "/analytics/statistics/all" "$http_code" "$(echo "$body" | jq -c '.')" "200" "PASS"
else
    echo -e "${RED}❌ FAIL${NC} (HTTP $http_code)"
    FAIL=$((FAIL + 1))
    add_result "Analytics" "All Statistics" "GET" "/analytics/statistics/all" "$http_code" "\"$body\"" "200" "FAIL"
fi
echo "$body" | jq '.' 2>/dev/null
echo ""

#==============================================================================
# FINALIZE JSON OUTPUT
#==============================================================================
echo "  ]," >> "$OUTPUT_FILE"
echo "  \"summary\": {" >> "$OUTPUT_FILE"
echo "    \"total\": $TOTAL," >> "$OUTPUT_FILE"
echo "    \"passed\": $PASS," >> "$OUTPUT_FILE"
echo "    \"failed\": $FAIL," >> "$OUTPUT_FILE"
echo "    \"success_rate\": \"$(awk "BEGIN {printf \"%.2f\", ($PASS/$TOTAL)*100}")%\"" >> "$OUTPUT_FILE"
echo "  }" >> "$OUTPUT_FILE"
echo "}" >> "$OUTPUT_FILE"

#==============================================================================
# SUMMARY
#==============================================================================
echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║                     TEST SUMMARY                               ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""
echo "Total Tests:  $TOTAL"
echo -e "${GREEN}Passed:       $PASS${NC}"
echo -e "${RED}Failed:       $FAIL${NC}"
SUCCESS_RATE=$(awk "BEGIN {printf \"%.2f\", ($PASS/$TOTAL)*100}")
echo "Success Rate: $SUCCESS_RATE%"
echo ""
echo "Results saved to: $OUTPUT_FILE"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║            🎉 ALL TESTS PASSED SUCCESSFULLY! 🎉                ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
    exit 0
else
    echo -e "${RED}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║              ⚠️  SOME TESTS FAILED ⚠️                          ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════════════════════════════╝${NC}"
    exit 1
fi
