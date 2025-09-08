package runner

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cassielabs/httpx/internal/executor"
	"github.com/cassielabs/httpx/internal/parser"
)

func RunTests(filePath string, timeout time.Duration) error {
	httpFile, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	for key, value := range httpFile.Variables {
		if envValue := os.Getenv(key); envValue != "" {
			httpFile.Variables[key] = envValue
		} else {
			httpFile.Variables[key] = value
		}
	}

	exec := executor.New(timeout)
	
	totalTests := len(httpFile.Requests)
	passed := 0
	failed := 0

	fmt.Printf("Running %d tests from %s\n\n", totalTests, filePath)

	for i, req := range httpFile.Requests {
		req.ApplyVariables(httpFile.Variables)
		
		testName := fmt.Sprintf("Test %d: %s %s", i+1, req.Method, req.URL)
		if req.Name != "" {
			testName = fmt.Sprintf("Test %d [%s]: %s %s", i+1, req.Name, req.Method, req.URL)
		}
		
		fmt.Printf("Running %s... ", testName)
		
		resp, err := exec.Execute(req)
		if err != nil {
			fmt.Printf("❌ FAILED\n")
			fmt.Printf("  Error: %v\n", err)
			failed++
			continue
		}
		
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			fmt.Printf("✅ PASSED (Status: %d, Duration: %v)\n", resp.StatusCode, resp.Duration)
			passed++
		} else {
			fmt.Printf("❌ FAILED\n")
			fmt.Printf("  Status: %d %s\n", resp.StatusCode, resp.Status)
			if resp.Body != "" && len(resp.Body) < 200 {
				fmt.Printf("  Body: %s\n", strings.TrimSpace(resp.Body))
			}
			failed++
		}
	}

	fmt.Print("\n" + strings.Repeat("-", 50) + "\n")
	fmt.Printf("Test Results: %d/%d passed", passed, totalTests)
	
	if failed > 0 {
		fmt.Printf(" (%d failed)\n", failed)
		return fmt.Errorf("%d tests failed", failed)
	}
	
	fmt.Printf(" ✅\n")
	return nil
}