package parser

import (
	"bufio"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	requestLineRegex = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS|TRACE|CONNECT)\s+(.+?)(?:\s+HTTP/[\d.]+)?$`)
	variableRegex    = regexp.MustCompile(`\{\{(.+?)\}\}`)
	separatorRegex   = regexp.MustCompile(`^###\s*(.*)$`)
)

func ParseFile(path string) (*HTTPFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	httpFile := &HTTPFile{
		Path:      path,
		Requests:  []HTTPRequest{},
		Variables: make(map[string]string),
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var currentRequest *HTTPRequest
	inBody := false
	bodyLines := []string{}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if separatorRegex.MatchString(line) {
			if currentRequest != nil && currentRequest.Method != "" {
				if len(bodyLines) > 0 {
					currentRequest.Body = strings.Join(bodyLines, "\n")
				}
				httpFile.Requests = append(httpFile.Requests, *currentRequest)
			}
			match := separatorRegex.FindStringSubmatch(line)
			currentRequest = &HTTPRequest{
				Headers: make(http.Header),
				Name:    strings.TrimSpace(match[1]),
			}
			inBody = false
			bodyLines = []string{}
			continue
		}

		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "@") {
			parts := strings.SplitN(line[1:], "=", 2)
			if len(parts) == 2 {
				httpFile.Variables[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
			continue
		}

		if currentRequest == nil {
			if strings.TrimSpace(line) == "" {
				continue
			}
			currentRequest = &HTTPRequest{
				Headers: make(http.Header),
			}
		}

		if inBody {
			bodyLines = append(bodyLines, line)
			continue
		}

		if requestLineRegex.MatchString(line) {
			matches := requestLineRegex.FindStringSubmatch(line)
			currentRequest.Method = matches[1]
			currentRequest.URL = matches[2]
			continue
		}

		if strings.Contains(line, ":") && currentRequest.Method != "" {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				headerName := strings.TrimSpace(parts[0])
				headerValue := strings.TrimSpace(parts[1])
				currentRequest.Headers.Add(headerName, headerValue)
			}
			continue
		}

		if strings.TrimSpace(line) == "" && currentRequest.Method != "" {
			inBody = true
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "http://") || strings.HasPrefix(strings.TrimSpace(line), "https://") {
			currentRequest.Method = "GET"
			currentRequest.URL = strings.TrimSpace(line)
			continue
		}
	}

	if currentRequest != nil && currentRequest.Method != "" {
		if len(bodyLines) > 0 {
			currentRequest.Body = strings.Join(bodyLines, "\n")
		}
		httpFile.Requests = append(httpFile.Requests, *currentRequest)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return httpFile, nil
}

func ReplaceVariables(text string, variables map[string]string) string {
	return variableRegex.ReplaceAllStringFunc(text, func(match string) string {
		varName := variableRegex.FindStringSubmatch(match)[1]
		if value, ok := variables[varName]; ok {
			return value
		}
		return match
	})
}

func (r *HTTPRequest) ApplyVariables(variables map[string]string) {
	r.URL = ReplaceVariables(r.URL, variables)
	r.Body = ReplaceVariables(r.Body, variables)
	
	for key, values := range r.Headers {
		for i, value := range values {
			r.Headers[key][i] = ReplaceVariables(value, variables)
		}
	}
}

func ParseString(content string) (*HTTPFile, error) {
	tmpFile, err := os.CreateTemp("", "*.http")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		return nil, err
	}
	tmpFile.Close()

	return ParseFile(tmpFile.Name())
}