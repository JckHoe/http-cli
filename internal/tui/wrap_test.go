package tui

import (
	"strings"
	"testing"
)

func TestWrapContent_ShortLine(t *testing.T) {
	content := "Hello World"
	width := 20

	result := wrapContent(content, width)

	if result != content {
		t.Errorf("Expected short line to remain unchanged, got %s", result)
	}
}

func TestWrapContent_ExactWidth(t *testing.T) {
	content := "1234567890"
	width := 10

	result := wrapContent(content, width)

	if result != content {
		t.Errorf("Expected line at exact width to remain unchanged, got %s", result)
	}
}

func TestWrapContent_LongLine(t *testing.T) {
	content := "This is a very long line that exceeds the width"
	width := 20

	result := wrapContent(content, width)
	lines := strings.Split(result, "\n")

	if len(lines) < 2 {
		t.Errorf("Expected line to be wrapped into multiple lines, got %d lines", len(lines))
	}

	for _, line := range lines {
		if len(line) > width {
			t.Errorf("Expected all lines to be <= %d chars, got line with %d chars: %s", width, len(line), line)
		}
	}

	unwrapped := strings.ReplaceAll(result, "\n", "")
	if unwrapped != content {
		t.Errorf("Content changed after wrapping. Expected: %s, Got: %s", content, unwrapped)
	}
}

func TestWrapContent_MultipleLines(t *testing.T) {
	content := "Line 1\nLine 2 is much longer than the width\nLine 3"
	width := 20

	result := wrapContent(content, width)
	lines := strings.Split(result, "\n")

	if len(lines) < 3 {
		t.Errorf("Expected at least 3 lines, got %d", len(lines))
	}

	for _, line := range lines {
		if len(line) > width {
			t.Errorf("Expected all lines to be <= %d chars, got line with %d chars: %s", width, len(line), line)
		}
	}
}

func TestWrapContent_JSONWithLongString(t *testing.T) {
	content := `{
  "jwt_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ5b3VyLWFwcGxpY2F0aW9uLW5hbWUiLCJzdWIiOiJhZG1pbkBlbWFpbC5jb20iLCJleHAiOjE3NjMxMDk2MTUsIm5iZiI6MTc2MzAyMzIxNSwiaWF0IjoxNzYzMDIzMjE1LCJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbkBlbWFpbC5jb20iLCJyb2xlIjoiYWRtaW4iLCJkaXNwbGF5X25hbWUiOiIiLCJmaXJzdF9uYW1lIjoiIiwibGFzdF9uYW1lIjoiIiwicGhvbmVfbnVtYmVyIjoiIiwiY29tcGFueSI6IiJ9.iHfPkJu6m2arh7tjKiOyBW9DomjaxPXvFxRV2Ktb__8"
}`
	width := 80

	result := wrapContent(content, width)
	lines := strings.Split(result, "\n")

	for _, line := range lines {
		if len(line) > width {
			t.Errorf("Expected all lines to be <= %d chars, got line with %d chars", width, len(line))
		}
	}

	unwrapped := strings.ReplaceAll(result, "\n", "")
	originalUnwrapped := strings.ReplaceAll(content, "\n", "")
	if unwrapped != originalUnwrapped {
		t.Errorf("Content was modified during wrapping")
	}
}

func TestWrapContent_PreservesAllCharacters(t *testing.T) {
	content := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+-=[]{}|;:',.<>?/`~"
	width := 20

	result := wrapContent(content, width)

	unwrapped := strings.ReplaceAll(result, "\n", "")
	if unwrapped != content {
		t.Errorf("Characters were lost or modified. Expected: %s, Got: %s", content, unwrapped)
	}
}

func TestWrapContent_EmptyString(t *testing.T) {
	content := ""
	width := 20

	result := wrapContent(content, width)

	if result != "" {
		t.Errorf("Expected empty string to remain empty, got %s", result)
	}
}

func TestWrapContent_ZeroWidth(t *testing.T) {
	content := "Hello World"
	width := 0

	result := wrapContent(content, width)

	if result != content {
		t.Errorf("Expected content to be unchanged with zero width, got %s", result)
	}
}

func TestWrapContent_NegativeWidth(t *testing.T) {
	content := "Hello World"
	width := -1

	result := wrapContent(content, width)

	if result != content {
		t.Errorf("Expected content to be unchanged with negative width, got %s", result)
	}
}

func TestWrapContent_NewlinesPreserved(t *testing.T) {
	content := "Line 1\nLine 2\nLine 3"
	width := 100

	result := wrapContent(content, width)

	if result != content {
		t.Errorf("Expected newlines to be preserved when no wrapping needed, got %s", result)
	}

	lineCount := strings.Count(result, "\n")
	expectedCount := strings.Count(content, "\n")
	if lineCount != expectedCount {
		t.Errorf("Expected %d newlines, got %d", expectedCount, lineCount)
	}
}

func TestWrapContent_UTF8Characters(t *testing.T) {
	content := "Hello 世界 🌍 Ñoño"
	width := 10

	result := wrapContent(content, width)
	lines := strings.Split(result, "\n")

	for _, line := range lines {
		if len(line) > width {
			t.Errorf("Expected all lines to be <= %d chars, got line with %d chars: %s", width, len(line), line)
		}
	}

	unwrapped := strings.ReplaceAll(result, "\n", "")
	if unwrapped != content {
		t.Errorf("UTF-8 characters were corrupted. Expected: %s, Got: %s", content, unwrapped)
	}
}

func TestWrapContent_TabsAndSpaces(t *testing.T) {
	content := "  Line with leading spaces\n\tLine with leading tab"
	width := 20

	result := wrapContent(content, width)

	if !strings.Contains(result, "  Line") {
		t.Errorf("Leading spaces were removed")
	}
	if !strings.Contains(result, "\tLine") {
		t.Errorf("Leading tabs were removed")
	}
}
