package executor

import (
	"testing"

	"github.com/cassielabs/hrun/internal/parser"
)

func TestApplyCaptureRules(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		captures []parser.CaptureRule
		expected map[string]string
	}{
		{
			name: "Simple JSON path",
			body: `{"token": "abc123", "userId": 42}`,
			captures: []parser.CaptureRule{
				{VariableName: "token", JSONPath: "token"},
				{VariableName: "userId", JSONPath: "userId"},
			},
			expected: map[string]string{
				"token":  "abc123",
				"userId": "42",
			},
		},
		{
			name: "Nested JSON path",
			body: `{"user": {"id": 1, "name": "John"}, "address": {"city": "NYC"}}`,
			captures: []parser.CaptureRule{
				{VariableName: "userId", JSONPath: "user.id"},
				{VariableName: "userName", JSONPath: "user.name"},
				{VariableName: "city", JSONPath: "address.city"},
			},
			expected: map[string]string{
				"userId":   "1",
				"userName": "John",
				"city":     "NYC",
			},
		},
		{
			name: "Array access",
			body: `{"items": [{"id": 1}, {"id": 2}]}`,
			captures: []parser.CaptureRule{
				{VariableName: "firstId", JSONPath: "items.0.id"},
				{VariableName: "secondId", JSONPath: "items.1.id"},
			},
			expected: map[string]string{
				"firstId":  "1",
				"secondId": "2",
			},
		},
		{
			name: "Missing key",
			body: `{"token": "abc123"}`,
			captures: []parser.CaptureRule{
				{VariableName: "token", JSONPath: "token"},
				{VariableName: "missing", JSONPath: "doesNotExist"},
			},
			expected: map[string]string{
				"token": "abc123",
			},
		},
		{
			name:     "Invalid JSON",
			body:     `not json`,
			captures: []parser.CaptureRule{
				{VariableName: "token", JSONPath: "token"},
			},
			expected: map[string]string{},
		},
		{
			name: "Empty captures",
			body: `{"token": "abc123"}`,
			captures: []parser.CaptureRule{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyCaptureRules(tt.body, tt.captures)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d captured variables, got %d", len(tt.expected), len(result))
			}

			for key, expectedValue := range tt.expected {
				if actualValue, ok := result[key]; !ok {
					t.Errorf("Expected key %q not found in result", key)
				} else if actualValue != expectedValue {
					t.Errorf("For key %q: expected %q, got %q", key, expectedValue, actualValue)
				}
			}

			for key := range result {
				if _, ok := tt.expected[key]; !ok {
					t.Errorf("Unexpected key %q in result", key)
				}
			}
		})
	}
}
