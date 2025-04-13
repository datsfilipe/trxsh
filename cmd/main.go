package main

import (
	"fmt"
	"os"

	"github.com/datsfilipe/trxsh/pkg/cli"
	"github.com/datsfilipe/trxsh/pkg/integrations"
	"github.com/datsfilipe/trxsh/pkg/register"
)

func printUsage() {
	fmt.Printf("Usage: %s [OPTIONS] [FILES]\n", os.Args[0])
	fmt.Println("Options:")
	fmt.Println("  --fzf, -f       : Restore files using fzf")
	fmt.Println("  --list, -l       : List files in trash")
	fmt.Println("  --restore, -r ID : Restore file by ID")
	fmt.Println("  --cleanup, -c    : Empty all trash directories")
	fmt.Println("  --help, -h       : Show this help")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	c, err := cli.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--fzf", "-f":
		reg, err := register.New("")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := reg.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fzf := integrations.NewFzf(reg)
		if err := fzf.RestoreWithFzf(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "--list", "-l":
		if err := c.List(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "--restore", "-r":
		if len(os.Args) < 3 {
			fmt.Println("Please provide an ID")
			os.Exit(1)
		}
		if err := c.Restore(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "--cleanup", "-c":
		if err := c.Cleanup(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Trash emptied")

	case "--help", "-h":
		printUsage()

	default:
		if err := c.Trash(os.Args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
