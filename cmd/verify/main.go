package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/verify"
)

func main() {
	host := flag.String("host", "", "Target server IP or hostname (not required with --local)")
	keyPath := flag.String("key", "", "Path to SSH private key (not required with --local)")
	user := flag.String("user", "root", "SSH username")
	port := flag.Int("port", 22, "SSH port")
	prebaked := flag.Bool("prebaked", false, "Run pre-baked image checks")
	deployed := flag.Bool("deployed", false, "Run deployed server checks")
	local := flag.Bool("local", false, "Run checks locally instead of via SSH")
	timeout := flag.Duration("timeout", 30*time.Second, "SSH connection timeout")

	flag.Parse()

	// Validate mode selection
	if !*prebaked && !*deployed {
		fmt.Fprintln(os.Stderr, "Error: either --prebaked or --deployed must be specified")
		flag.Usage()
		os.Exit(1)
	}

	if *prebaked && *deployed {
		fmt.Fprintln(os.Stderr, "Error: only one of --prebaked or --deployed can be specified")
		flag.Usage()
		os.Exit(1)
	}

	// Validate required flags for SSH mode
	if !*local {
		if *host == "" {
			fmt.Fprintln(os.Stderr, "Error: --host is required (or use --local for local checks)")
			flag.Usage()
			os.Exit(1)
		}

		if *keyPath == "" {
			fmt.Fprintln(os.Stderr, "Error: --key is required (or use --local for local checks)")
			flag.Usage()
			os.Exit(1)
		}
	}

	// Determine mode and checks
	var mode string
	var checks []verify.Check

	if *prebaked {
		mode = "prebaked"
		checks = verify.PrebakedChecks()
	} else {
		mode = "deployed"
		checks = verify.DeployedChecks()
	}

	var results []verify.Result
	var displayHost string

	if *local {
		// Run checks locally
		displayHost = "localhost"
		fmt.Println("Running local verification...")
		results = verify.RunChecksLocal(checks)
	} else {
		// Connect to the server via SSH
		displayHost = *host
		config := verify.SSHConfig{
			Host:    *host,
			Port:    *port,
			User:    *user,
			KeyPath: *keyPath,
			Timeout: *timeout,
		}

		fmt.Printf("Connecting to %s@%s:%d...\n", config.User, config.Host, config.Port)

		client, err := verify.Connect(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		results = verify.RunChecks(client, checks)
	}

	// Print results and exit with appropriate code
	exitCode := verify.PrintResults(displayHost, mode, results)
	os.Exit(exitCode)
}
