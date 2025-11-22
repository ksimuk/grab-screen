package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/ksimuk/grab-screen/portal"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	keepFlag := flag.Bool("k", false, "Keep the screenshot file after execution")
	flag.BoolVar(keepFlag, "keep", false, "Keep the screenshot file after execution")

	versionFlag := flag.Bool("v", false, "Print version information")
	flag.BoolVar(versionFlag, "version", false, "Print version information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [command [args...]]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIf a command is provided, the screenshot path will be appended to its arguments.\n")
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s swappy -f\n", os.Args[0])
	}

	flag.Parse()

	if *versionFlag {
		fmt.Printf("grab-screen version %s\n", version)
		return nil
	}

	// Take the screenshot
	screenshotPath, err := portal.TakeScreenshot()
	if err != nil {
		return fmt.Errorf("error taking screenshot: %w", err)
	}

	if screenshotPath == "" {
		return fmt.Errorf("error: received empty screenshot path")
	}

	// Ensure cleanup unless keep flag is set
	if !*keepFlag {
		defer func() {
			if err := os.Remove(screenshotPath); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to remove temporary file %s: %v\n", screenshotPath, err)
			}
		}()
	}

	args := flag.Args()
	if len(args) > 0 {
		// Prepare command
		cmdName := args[0]
		cmdArgs := args[1:]
		cmdArgs = append(cmdArgs, screenshotPath)

		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			// If the command fails, we return the error, and defer will handle cleanup
			return fmt.Errorf("command execution failed: %w", err)
		}
	} else {
		// If no command provided, just print the path
		fmt.Println(screenshotPath)
	}

	return nil
}
