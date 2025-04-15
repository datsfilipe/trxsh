package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/datsfilipe/trxsh/pkg/cli"
	"github.com/datsfilipe/trxsh/pkg/integrations"
)

func printUsage() {
	fmt.Printf("Usage: %s [OPTIONS] [FILES]\n", os.Args[0])
	fmt.Println("Options:")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "  --fzf, -f\t: Restore files using fzf")
	fmt.Fprintln(w, "  --list, -l\t: List files in trash")
	fmt.Fprintln(w, "  --restore, -r ID\t: Restore file by ID")
	fmt.Fprintln(w, "  --cleanup, -c\t: Empty all trash directories")
	fmt.Fprintln(w, "  --dir-sizes, -s\t: Show directory sizes")
	fmt.Fprintln(w, "  --help, -h\t: Show this help")
	w.Flush()
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
		fzf, err := integrations.NewFzf()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

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

	case "--dir-sizes", "-s":
		err := c.PrintDirSizes()
		if err != nil {
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
