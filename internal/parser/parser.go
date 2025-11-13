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
	captureRegex     = regexp.MustCompile(`^@capture\s+(\w+)\s*=\s*(.+)$`)
)

func ParseFile(path string) (*HTTPFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

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
	descriptionLines := []string{}

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
				Headers:    make(http.Header),
				Name:       strings.TrimSpace(match[1]),
				LineNumber: lineNum,
			}
			inBody = false
			bodyLines = []string{}
			descriptionLines = []string{}
			continue
		}

		if inBody {
			bodyLines = append(bodyLines, line)
			continue
		}

		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			if currentRequest != nil && currentRequest.Method == "" {
				comment := strings.TrimSpace(line)
				if strings.HasPrefix(comment, "#") {
					comment = strings.TrimSpace(strings.TrimPrefix(comment, "#"))
				} else if strings.HasPrefix(comment, "//") {
					comment = strings.TrimSpace(strings.TrimPrefix(comment, "//"))
				}
				if comment != "" {
					if matches := captureRegex.FindStringSubmatch(comment); len(matches) == 3 {
						currentRequest.Captures = append(currentRequest.Captures, CaptureRule{
							VariableName: matches[1],
							JSONPath:     strings.TrimSpace(matches[2]),
						})
					} else {
						descriptionLines = append(descriptionLines, comment)
					}
				}
			}
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

		if requestLineRegex.MatchString(line) {
			matches := requestLineRegex.FindStringSubmatch(line)
			currentRequest.Method = matches[1]
			currentRequest.URL = matches[2]
			if len(descriptionLines) > 0 {
				currentRequest.Description = strings.Join(descriptionLines, " ")
				descriptionLines = []string{}
			}
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "HTTP/") {
			continue
		}

		trimmedLine := strings.TrimSpace(line)
		if strings.Contains(line, ":") && currentRequest.Method != "" {
			if strings.HasPrefix(trimmedLine, "{") || strings.HasPrefix(trimmedLine, "[") || strings.HasPrefix(trimmedLine, "<") {
				inBody = true
				bodyLines = append(bodyLines, line)
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				headerName := strings.TrimSpace(parts[0])
				headerValue := strings.TrimSpace(parts[1])
				if headerName != "" && !strings.ContainsAny(headerName, " \t\"{}<>[]") {
					currentRequest.Headers.Add(headerName, headerValue)
				} else {
					inBody = true
					bodyLines = append(bodyLines, line)
				}
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
			if len(descriptionLines) > 0 {
				currentRequest.Description = strings.Join(descriptionLines, " ")
				descriptionLines = []string{}
			}
			continue
		}

		if currentRequest.Method != "" {
			inBody = true
			bodyLines = append(bodyLines, line)
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
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err := tmpFile.WriteString(content); err != nil {
		return nil, err
	}
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}

	return ParseFile(tmpFile.Name())
}