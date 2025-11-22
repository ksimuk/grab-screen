package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/ksimuk/grab-screen/portal"
)

func main() {
	keepFlag := flag.Bool("k", false, "Keep the screenshot file after execution")
	flag.BoolVar(keepFlag, "keep", false, "Keep the screenshot file after execution")
	flag.Parse()

	// Take the screenshot
	screenshotPath, err := portal.TakeScreenshot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error taking screenshot: %v\n", err)
		os.Exit(1)
	}

	if screenshotPath == "" {
		fmt.Fprintf(os.Stderr, "Error: received empty screenshot path\n")
		os.Exit(1)
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
			// If the command fails, we still exit with error, but defer will handle cleanup
			fmt.Fprintf(os.Stderr, "Command execution failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		// If no command provided, just print the path
		fmt.Println(screenshotPath)
	}
}
