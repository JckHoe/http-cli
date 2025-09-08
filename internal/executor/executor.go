package executor

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cassielabs/httpx/internal/parser"
)

type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Body       string
	Duration   time.Duration
	Error      error
}

type Executor struct {
	client  *http.Client
	timeout time.Duration
}

func New(timeout time.Duration) *Executor {
	return &Executor{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

func (e *Executor) Execute(req parser.HTTPRequest) (*Response, error) {
	start := time.Now()
	
	httpReq, err := http.NewRequest(req.Method, req.URL, strings.NewReader(req.Body))
	if err != nil {
		return &Response{
			Error:    err,
			Duration: time.Since(start),
		}, err
	}

	for key, values := range req.Headers {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}

	if req.Body != "" && httpReq.Header.Get("Content-Type") == "" {
		if strings.HasPrefix(strings.TrimSpace(req.Body), "{") || strings.HasPrefix(strings.TrimSpace(req.Body), "[") {
			httpReq.Header.Set("Content-Type", "application/json")
		} else if strings.HasPrefix(strings.TrimSpace(req.Body), "<") {
			httpReq.Header.Set("Content-Type", "application/xml")
		} else {
			httpReq.Header.Set("Content-Type", "text/plain")
		}
	}

	resp, err := e.client.Do(httpReq)
	if err != nil {
		return &Response{
			Error:    err,
			Duration: time.Since(start),
		}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Headers:    resp.Header,
			Error:      err,
			Duration:   time.Since(start),
		}, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Body:       string(body),
		Duration:   time.Since(start),
	}, nil
}

func (e *Executor) ExecuteAll(file *parser.HTTPFile) ([]*Response, error) {
	responses := make([]*Response, 0, len(file.Requests))
	
	for _, req := range file.Requests {
		req.ApplyVariables(file.Variables)
		resp, err := e.Execute(req)
		if err != nil && resp == nil {
			return responses, err
		}
		responses = append(responses, resp)
	}
	
	return responses, nil
}

func FormatResponse(resp *Response) string {
	var buf bytes.Buffer
	
	if resp.Error != nil {
		fmt.Fprintf(&buf, "Error: %v\n", resp.Error)
		fmt.Fprintf(&buf, "Duration: %v\n", resp.Duration)
		return buf.String()
	}
	
	fmt.Fprintf(&buf, "Status: %s\n", resp.Status)
	fmt.Fprintf(&buf, "Duration: %v\n", resp.Duration)
	fmt.Fprintln(&buf, "\nHeaders:")
	for key, values := range resp.Headers {
		for _, value := range values {
			fmt.Fprintf(&buf, "  %s: %s\n", key, value)
		}
	}
	
	if resp.Body != "" {
		fmt.Fprintln(&buf, "\nBody:")
		fmt.Fprintln(&buf, resp.Body)
	}
	
	return buf.String()
}