package parser

import (
	"testing"
)

func TestParseCaptureDirective(t *testing.T) {
	content := `### Login Request
# @capture token=$.jwt_token
# @capture userId=$.user.id
POST https://api.example.com/login
Content-Type: application/json

{
  "username": "test",
  "password": "pass"
}

###

### Test nested capture
# @capture city=$.address.city
# @capture company=$.company.name
GET https://api.example.com/users/1
`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	if len(httpFile.Requests) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(httpFile.Requests))
	}

	firstReq := httpFile.Requests[0]
	if len(firstReq.Captures) != 2 {
		t.Fatalf("Expected 2 captures in first request, got %d", len(firstReq.Captures))
	}

	if firstReq.Captures[0].VariableName != "token" || firstReq.Captures[0].JSONPath != "$.jwt_token" {
		t.Errorf("First capture incorrect: got %+v", firstReq.Captures[0])
	}

	if firstReq.Captures[1].VariableName != "userId" || firstReq.Captures[1].JSONPath != "$.user.id" {
		t.Errorf("Second capture incorrect: got %+v", firstReq.Captures[1])
	}

	secondReq := httpFile.Requests[1]
	if len(secondReq.Captures) != 2 {
		t.Fatalf("Expected 2 captures in second request, got %d", len(secondReq.Captures))
	}

	if secondReq.Captures[0].VariableName != "city" || secondReq.Captures[0].JSONPath != "$.address.city" {
		t.Errorf("Third capture incorrect: got %+v", secondReq.Captures[0])
	}
}

func TestCaptureWithComments(t *testing.T) {
	content := `### Test Request
# This is a regular comment
# @capture token=$.jwt_token
# Another regular comment
POST https://api.example.com/login
`

	httpFile, err := ParseString(content)
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	if len(httpFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(httpFile.Requests))
	}

	req := httpFile.Requests[0]
	if len(req.Captures) != 1 {
		t.Fatalf("Expected 1 capture, got %d", len(req.Captures))
	}

	if req.Description != "This is a regular comment Another regular comment" {
		t.Errorf("Description incorrect: got %q", req.Description)
	}
}
