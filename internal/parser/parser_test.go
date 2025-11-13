package parser

import (
	"strings"
	"testing"
)

func TestParseFile_SimpleRequest(t *testing.T) {
	content := `### Get Users
GET https://api.example.com/users`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(httpFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(httpFile.Requests))
	}

	req := httpFile.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Expected method GET, got %s", req.Method)
	}
	if req.URL != "https://api.example.com/users" {
		t.Errorf("Expected URL https://api.example.com/users, got %s", req.URL)
	}
	if req.Name != "Get Users" {
		t.Errorf("Expected name 'Get Users', got %s", req.Name)
	}
}

func TestParseFile_RequestWithHeaders(t *testing.T) {
	content := `### Create User
POST https://api.example.com/users
Content-Type: application/json
Authorization: Bearer token123`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(httpFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(httpFile.Requests))
	}

	req := httpFile.Requests[0]
	if req.Method != "POST" {
		t.Errorf("Expected method POST, got %s", req.Method)
	}

	if len(req.Headers) != 2 {
		t.Fatalf("Expected 2 headers, got %d", len(req.Headers))
	}

	if req.Headers.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header, got %s", req.Headers.Get("Content-Type"))
	}

	if req.Headers.Get("Authorization") != "Bearer token123" {
		t.Errorf("Expected Authorization header, got %s", req.Headers.Get("Authorization"))
	}
}

func TestParseFile_RequestWithBody(t *testing.T) {
	content := `### Create User
POST https://api.example.com/users
Content-Type: application/json

{
  "username": "test@example.com",
  "password": "secret"
}`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(httpFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(httpFile.Requests))
	}

	req := httpFile.Requests[0]
	expectedBody := `{
  "username": "test@example.com",
  "password": "secret"
}`

	if strings.TrimSpace(req.Body) != strings.TrimSpace(expectedBody) {
		t.Errorf("Expected body:\n%s\nGot:\n%s", expectedBody, req.Body)
	}
}

func TestParseFile_RequestWithBodyNoBlankLine(t *testing.T) {
	content := `### Admin Login
POST http://localhost:8080/parking/v1/auth/login
content-type: application/json
{
    "username": "admin@email.com",
    "password": "password"
}`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(httpFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(httpFile.Requests))
	}

	req := httpFile.Requests[0]
	if req.Method != "POST" {
		t.Errorf("Expected method POST, got %s", req.Method)
	}

	if req.Headers.Get("content-type") != "application/json" {
		t.Errorf("Expected content-type header, got %s", req.Headers.Get("content-type"))
	}

	if !strings.Contains(req.Body, `"username": "admin@email.com"`) {
		t.Errorf("Expected body to contain username field, got: %s", req.Body)
	}
}

func TestParseFile_RequestWithHTTPVersion(t *testing.T) {
	content := `### Health Check
GET http://localhost:8080/health HTTP/1.1`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(httpFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(httpFile.Requests))
	}

	req := httpFile.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Expected method GET, got %s", req.Method)
	}
	if req.URL != "http://localhost:8080/health" {
		t.Errorf("Expected URL without HTTP/1.1, got %s", req.URL)
	}
}

func TestParseFile_RequestWithSeparateHTTPVersion(t *testing.T) {
	content := `### Health Check
GET http://localhost:8080/health
HTTP/1.1`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(httpFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(httpFile.Requests))
	}

	req := httpFile.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Expected method GET, got %s", req.Method)
	}
}

func TestParseFile_MultipleRequests(t *testing.T) {
	content := `### Get Users
GET https://api.example.com/users

### Create User
POST https://api.example.com/users
Content-Type: application/json

{
  "username": "test"
}

### Delete User
DELETE https://api.example.com/users/1`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(httpFile.Requests) != 3 {
		t.Fatalf("Expected 3 requests, got %d", len(httpFile.Requests))
	}

	if httpFile.Requests[0].Method != "GET" {
		t.Errorf("Expected first request to be GET, got %s", httpFile.Requests[0].Method)
	}

	if httpFile.Requests[1].Method != "POST" {
		t.Errorf("Expected second request to be POST, got %s", httpFile.Requests[1].Method)
	}

	if httpFile.Requests[2].Method != "DELETE" {
		t.Errorf("Expected third request to be DELETE, got %s", httpFile.Requests[2].Method)
	}
}

func TestParseFile_Variables(t *testing.T) {
	content := `@baseUrl = https://api.example.com
@token = abc123

### Get Users
GET {{baseUrl}}/users
Authorization: Bearer {{token}}`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(httpFile.Variables) != 2 {
		t.Fatalf("Expected 2 variables, got %d", len(httpFile.Variables))
	}

	if httpFile.Variables["baseUrl"] != "https://api.example.com" {
		t.Errorf("Expected baseUrl variable, got %s", httpFile.Variables["baseUrl"])
	}

	if httpFile.Variables["token"] != "abc123" {
		t.Errorf("Expected token variable, got %s", httpFile.Variables["token"])
	}
}

func TestParseFile_VariableReplacement(t *testing.T) {
	content := `@baseUrl = https://api.example.com
@token = abc123

### Get Users
GET {{baseUrl}}/users
Authorization: Bearer {{token}}`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	req := httpFile.Requests[0]
	req.ApplyVariables(httpFile.Variables)

	if req.URL != "https://api.example.com/users" {
		t.Errorf("Expected URL with variable replaced, got %s", req.URL)
	}

	if req.Headers.Get("Authorization") != "Bearer abc123" {
		t.Errorf("Expected Authorization with variable replaced, got %s", req.Headers.Get("Authorization"))
	}
}

func TestParseFile_Comments(t *testing.T) {
	content := `### Get Users
# This endpoint returns all users
GET https://api.example.com/users`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(httpFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(httpFile.Requests))
	}

	req := httpFile.Requests[0]
	if req.Description != "This endpoint returns all users" {
		t.Errorf("Expected description 'This endpoint returns all users', got %s", req.Description)
	}
}

func TestParseFile_EmptyFile(t *testing.T) {
	content := ``

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(httpFile.Requests) != 0 {
		t.Errorf("Expected 0 requests, got %d", len(httpFile.Requests))
	}
}

func TestParseFile_InvalidHeaderNotParsedAsHeader(t *testing.T) {
	content := `### Test
POST http://localhost:8080/api
content-type: application/json
{
    "key": "value"
}`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	req := httpFile.Requests[0]

	for headerName := range req.Headers {
		if strings.Contains(headerName, "{") || strings.Contains(headerName, `"`) {
			t.Errorf("JSON body was incorrectly parsed as header: %s", headerName)
		}
	}

	if !strings.Contains(req.Body, `"key": "value"`) {
		t.Errorf("Expected body to contain JSON, got: %s", req.Body)
	}
}
