package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cassielabs/httpx/internal/executor"
	"github.com/cassielabs/httpx/internal/parser"
	"github.com/cassielabs/httpx/internal/runner"
	"github.com/cassielabs/httpx/internal/tui"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	requestIndex int
	requestName  string
	envFile      string
	timeout      time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "httpx",
	Short: "HTTP file runner CLI tool",
	Long:  "A CLI tool to read and execute .http files with TUI support",
}

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Run HTTP requests from a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if envFile != "" {
			if err := godotenv.Load(envFile); err != nil {
				fmt.Printf("Warning: Could not load env file %s: %v\n", envFile, err)
			}
		}

		httpFile, err := parser.ParseFile(args[0])
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

		if requestName != "" {
			for _, req := range httpFile.Requests {
				if req.Name == requestName {
					req.ApplyVariables(httpFile.Variables)
					resp, err := exec.Execute(req)
					if err != nil {
						return err
					}
					fmt.Println(executor.FormatResponse(resp))
					return nil
				}
			}
			return fmt.Errorf("request with name '%s' not found", requestName)
		}

		if requestIndex > 0 {
			if requestIndex > len(httpFile.Requests) {
				return fmt.Errorf("request index %d out of range (file has %d requests)", requestIndex, len(httpFile.Requests))
			}
			req := httpFile.Requests[requestIndex-1]
			req.ApplyVariables(httpFile.Variables)
			resp, err := exec.Execute(req)
			if err != nil {
				return err
			}
			fmt.Println(executor.FormatResponse(resp))
			return nil
		}

		responses, err := exec.ExecuteAll(httpFile)
		if err != nil {
			return err
		}

		for i, resp := range responses {
			req := httpFile.Requests[i]
			fmt.Printf("\n=== Request %d: %s %s ===\n", i+1, req.Method, req.URL)
			if req.Name != "" {
				fmt.Printf("Name: %s\n", req.Name)
			}
			fmt.Println(executor.FormatResponse(resp))
		}

		return nil
	},
}

var tuiCmd = &cobra.Command{
	Use:   "tui [file]",
	Short: "Open file in TUI mode",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if envFile != "" {
			if err := godotenv.Load(envFile); err != nil {
				fmt.Printf("Warning: Could not load env file %s: %v\n", envFile, err)
			}
		}

		var filePath string
		if len(args) > 0 {
			filePath = args[0]
		}

		return tui.Run(filePath, timeout)
	},
}

var testCmd = &cobra.Command{
	Use:   "test [file]",
	Short: "Run tests for HTTP requests",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if envFile != "" {
			if err := godotenv.Load(envFile); err != nil {
				fmt.Printf("Warning: Could not load env file %s: %v\n", envFile, err)
			}
		}

		return runner.RunTests(args[0], timeout)
	},
}

func init() {
	runCmd.Flags().IntVar(&requestIndex, "request", 0, "Run specific request by index (1-based)")
	runCmd.Flags().StringVar(&requestName, "name", "", "Run specific request by name")
	runCmd.Flags().StringVar(&envFile, "env", "", "Environment file to load")
	runCmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")

	tuiCmd.Flags().StringVar(&envFile, "env", "", "Environment file to load")
	tuiCmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")

	testCmd.Flags().StringVar(&envFile, "env", "", "Environment file to load")
	testCmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(testCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}