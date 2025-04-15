package integrations

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/datsfilipe/trxsh/pkg/cli"
)

type FzfIntegration struct {
	cli *cli.CLI
}

func NewFzf() (*FzfIntegration, error) {
	c, err := cli.New()
	if err != nil {
		return nil, err
	}
	return &FzfIntegration{cli: c}, nil
}

func (f *FzfIntegration) RestoreWithFzf() error {
	_, err := exec.LookPath("fzf")
	if err != nil {
		return fmt.Errorf("fzf not installed")
	}

	records := f.cli.RawList()
	if len(records) == 0 {
		fmt.Println("No files in trash")
		return nil
	}

	tmpfile, err := os.CreateTemp("", "trash-fzf-")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	for _, record := range records {
		fmt.Fprintf(tmpfile, "%d: %s -> %s\n", record.ID, record.Name, record.Path)
	}
	tmpfile.Close()

	cmd := exec.Command("fzf", "--height", "40%", "--border")
	fh, err := os.Open(tmpfile.Name())
	if err != nil {
		return err
	}
	defer fh.Close()

	cmd.Stdin = fh
	var outBuf strings.Builder
	mw := io.MultiWriter(os.Stdout, &outBuf)
	cmd.Stdout = mw
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return nil
		}
		return err
	}

	selection := strings.TrimSpace(outBuf.String())
	if selection == "" {
		return nil
	}

	idStr := strings.Split(selection, ":")[0]
	return f.cli.Restore(idStr)
}
