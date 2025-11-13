package executor

import (
	"net/http"
	"strings"
	"testing"
)

func TestFormatBody_JSON(t *testing.T) {
	body := `{"username":"test","password":"secret"}`
	contentType := "application/json"

	result := formatBody(body, contentType)

	if !strings.Contains(result, `"username": "test"`) {
		t.Errorf("Expected pretty-printed JSON with proper spacing")
	}

	if !strings.Contains(result, "\n") {
		t.Errorf("Expected JSON to be formatted with newlines")
	}
}

func TestFormatBody_JSONWithoutContentType(t *testing.T) {
	body := `{"username":"test","password":"secret"}`
	contentType := ""

	result := formatBody(body, contentType)

	if !strings.Contains(result, `"username": "test"`) {
		t.Errorf("Expected JSON to be detected by content, got: %s", result)
	}
}

func TestFormatBody_JSONArray(t *testing.T) {
	body := `[{"id":1},{"id":2}]`
	contentType := "application/json"

	result := formatBody(body, contentType)

	if !strings.Contains(result, `"id": 1`) {
		t.Errorf("Expected pretty-printed JSON array")
	}
}

func TestFormatBody_InvalidJSON(t *testing.T) {
	body := `{"invalid json`
	contentType := "application/json"

	result := formatBody(body, contentType)

	if result != body {
		t.Errorf("Expected invalid JSON to be returned as-is, got: %s", result)
	}
}

func TestFormatBody_PlainText(t *testing.T) {
	body := "Hello World"
	contentType := "text/plain"

	result := formatBody(body, contentType)

	if result != body {
		t.Errorf("Expected plain text to be returned as-is, got: %s", result)
	}
}

func TestFormatBody_HTML(t *testing.T) {
	body := "<html><body>Hello</body></html>"
	contentType := "text/html"

	result := formatBody(body, contentType)

	if result != body {
		t.Errorf("Expected HTML to be returned as-is, got: %s", result)
	}
}

func TestFormatBody_EmptyBody(t *testing.T) {
	body := ""
	contentType := "application/json"

	result := formatBody(body, contentType)

	if result != body {
		t.Errorf("Expected empty body to be returned as-is")
	}
}

func TestFormatBody_NestedJSON(t *testing.T) {
	body := `{"user":{"name":"John","address":{"city":"NYC"}}}`
	contentType := "application/json"

	result := formatBody(body, contentType)

	if !strings.Contains(result, `"name": "John"`) {
		t.Errorf("Expected nested JSON to be pretty-printed")
	}

	if !strings.Contains(result, `"city": "NYC"`) {
		t.Errorf("Expected deeply nested JSON to be pretty-printed")
	}
}

func TestFormatBody_JSONWithNumbers(t *testing.T) {
	body := `{"age":30,"price":19.99,"count":0}`
	contentType := "application/json"

	result := formatBody(body, contentType)

	if !strings.Contains(result, `"age": 30`) {
		t.Errorf("Expected integer to be formatted correctly")
	}

	if !strings.Contains(result, `"price": 19.99`) {
		t.Errorf("Expected float to be formatted correctly")
	}
}

func TestFormatBody_JSONWithBoolean(t *testing.T) {
	body := `{"active":true,"deleted":false}`
	contentType := "application/json"

	result := formatBody(body, contentType)

	if !strings.Contains(result, `"active": true`) {
		t.Errorf("Expected boolean true to be formatted correctly")
	}

	if !strings.Contains(result, `"deleted": false`) {
		t.Errorf("Expected boolean false to be formatted correctly")
	}
}

func TestFormatBody_JSONWithNull(t *testing.T) {
	body := `{"value":null}`
	contentType := "application/json"

	result := formatBody(body, contentType)

	if !strings.Contains(result, `"value": null`) {
		t.Errorf("Expected null to be formatted correctly")
	}
}

func TestFormatResponse_Success(t *testing.T) {
	resp := &Response{
		StatusCode: 200,
		Status:     "200 OK",
		Headers:    http.Header{"Content-Type": []string{"application/json"}},
		Body:       `{"message":"success"}`,
	}

	result := FormatResponse(resp)

	if !strings.Contains(result, "Status: 200 OK") {
		t.Errorf("Expected status to be included in response")
	}

	if !strings.Contains(result, "Headers:") {
		t.Errorf("Expected headers section in response")
	}

	if !strings.Contains(result, "Body:") {
		t.Errorf("Expected body section in response")
	}

	if !strings.Contains(result, `"message": "success"`) {
		t.Errorf("Expected body to be pretty-printed JSON")
	}
}

func TestFormatResponse_Error(t *testing.T) {
	resp := &Response{
		Error: &testError{msg: "connection failed"},
	}

	result := FormatResponse(resp)

	if !strings.Contains(result, "Error: connection failed") {
		t.Errorf("Expected error message in response, got: %s", result)
	}

	if !strings.Contains(result, "Duration:") {
		t.Errorf("Expected duration in error response")
	}
}

func TestFormatResponse_WithMultipleHeaders(t *testing.T) {
	resp := &Response{
		StatusCode: 200,
		Status:     "200 OK",
		Headers: http.Header{
			"Content-Type":   []string{"application/json"},
			"Authorization":  []string{"Bearer token"},
			"Cache-Control":  []string{"no-cache"},
		},
		Body: "",
	}

	result := FormatResponse(resp)

	if !strings.Contains(result, "Content-Type: application/json") {
		t.Errorf("Expected Content-Type header in response")
	}

	if !strings.Contains(result, "Authorization: Bearer token") {
		t.Errorf("Expected Authorization header in response")
	}

	if !strings.Contains(result, "Cache-Control: no-cache") {
		t.Errorf("Expected Cache-Control header in response")
	}
}

func TestFormatResponse_EmptyBody(t *testing.T) {
	resp := &Response{
		StatusCode: 204,
		Status:     "204 No Content",
		Headers:    http.Header{},
		Body:       "",
	}

	result := FormatResponse(resp)

	if !strings.Contains(result, "Status: 204 No Content") {
		t.Errorf("Expected status in response")
	}

	bodyIndex := strings.Index(result, "Body:")
	if bodyIndex != -1 {
		afterBody := result[bodyIndex+5:]
		if strings.TrimSpace(afterBody) != "" {
			t.Errorf("Expected no body content for empty body")
		}
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
