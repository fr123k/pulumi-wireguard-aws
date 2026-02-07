package verify

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
)

// Check represents a single verification check to run
type Check struct {
	Name     string
	Command  string
	Expected string // "exit:0" for exit code check, or substring to match in output
	Equals   string // Exact match (trimmed output must equal this value)
}

// Result represents the outcome of a single check
type Result struct {
	Check  Check
	Passed bool
	Output string
	Error  error
}

// RunChecks executes all checks against the SSH client and returns results
func RunChecks(client *ssh.Client, checks []Check) []Result {
	results := make([]Result, 0, len(checks))

	for _, check := range checks {
		result := runCheck(client, check)
		results = append(results, result)
	}

	return results
}

// RunChecksLocal executes all checks locally and returns results
func RunChecksLocal(checks []Check) []Result {
	results := make([]Result, 0, len(checks))

	for _, check := range checks {
		result := runCheckLocal(check)
		results = append(results, result)
	}

	return results
}

// runCheckLocal executes a single check locally and returns the result
func runCheckLocal(check Check) Result {
	output, exitCode, err := executeCommandLocal(check.Command)

	result := Result{
		Check:  check,
		Output: output,
		Error:  err,
	}

	if err != nil {
		result.Passed = false
		return result
	}

	// Check based on expected value
	if check.Equals != "" {
		// Exact match (trimmed)
		result.Passed = strings.TrimSpace(output) == check.Equals
	} else if check.Expected == "exit:0" {
		result.Passed = exitCode == 0
	} else if strings.HasPrefix(check.Expected, "contains:") {
		expected := strings.TrimPrefix(check.Expected, "contains:")
		result.Passed = strings.Contains(output, expected)
	} else {
		// Default: substring match in output
		result.Passed = strings.Contains(output, check.Expected)
	}

	return result
}

// executeCommandLocal runs a command locally and returns output, exit code, and error
func executeCommandLocal(command string) (string, int, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		// Check if it's an exit error to get the exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return outputStr, status.ExitStatus(), nil
			}
		}
		return outputStr, -1, err
	}

	return outputStr, 0, nil
}

// runCheck executes a single check and returns the result
func runCheck(client *ssh.Client, check Check) Result {
	output, exitCode, err := executeCommand(client, check.Command)

	result := Result{
		Check:  check,
		Output: output,
		Error:  err,
	}

	if err != nil {
		result.Passed = false
		return result
	}

	// Check based on expected value
	if check.Equals != "" {
		// Exact match (trimmed)
		result.Passed = strings.TrimSpace(output) == check.Equals
	} else if check.Expected == "exit:0" {
		result.Passed = exitCode == 0
	} else if strings.HasPrefix(check.Expected, "contains:") {
		expected := strings.TrimPrefix(check.Expected, "contains:")
		result.Passed = strings.Contains(output, expected)
	} else {
		// Default: substring match in output
		result.Passed = strings.Contains(output, check.Expected)
	}

	return result
}

// executeCommand runs a command over SSH and returns output, exit code, and error
func executeCommand(client *ssh.Client, command string) (string, int, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", -1, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	outputStr := string(output)

	if err != nil {
		// Check if it's an exit error to get the exit code
		if exitErr, ok := err.(*ssh.ExitError); ok {
			return outputStr, exitErr.ExitStatus(), nil
		}
		return outputStr, -1, err
	}

	return outputStr, 0, nil
}

// PrintResults displays the check results in a formatted way
func PrintResults(host string, mode string, results []Result) int {
	fmt.Printf("\nVerifying server: %s (%s mode)\n\n", host, mode)

	passed := 0
	failed := 0

	for _, result := range results {
		if result.Passed {
			fmt.Printf("[✓] %s\n", result.Check.Name)
			passed++
		} else {
			failMsg := ""
			if result.Error != nil {
				failMsg = result.Error.Error()
			} else if result.Output != "" {
				// Trim and truncate output for display
				output := strings.TrimSpace(result.Output)
				if len(output) > 80 {
					output = output[:80] + "..."
				}
				failMsg = output
			} else {
				failMsg = "check failed"
			}
			fmt.Printf("[✗] %s - FAILED: %s\n", result.Check.Name, failMsg)
			failed++
		}
	}

	fmt.Printf("\nResults: %d/%d checks passed\n", passed, len(results))

	if failed > 0 {
		return 1
	}
	return 0
}
