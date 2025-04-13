package integrations

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/datsfilipe/trxsh/pkg/register"
)

type FzfIntegration struct {
	reg *register.Register
}

func NewFzf(reg *register.Register) *FzfIntegration {
	return &FzfIntegration{reg: reg}
}

func (f *FzfIntegration) RestoreWithFzf() error {
	_, err := exec.LookPath("fzf")
	if err != nil {
		return fmt.Errorf("fzf not installed")
	}
	records := f.reg.List()
	if len(records) == 0 {
		return fmt.Errorf("no files in trash")
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
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid selection: %s", selection)
	}
	return f.restoreFile(id)
}

func (f *FzfIntegration) restoreFile(id int) error {
	record, err := f.reg.Get(id)
	if err != nil {
		return err
	}
	trashRoot, err := register.GetTrashRoot(record.Path)
	if err != nil {
		return err
	}
	encodedName, err := register.EncodePath(record.Path)
	if err != nil {
		return err
	}
	trashPath := filepath.Join(trashRoot, encodedName)
	targetDir := filepath.Dir(record.Path)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}
	if err := os.Rename(trashPath, record.Path); err != nil {
		return err
	}
	if err := f.reg.Remove(id); err != nil {
		return err
	}
	fmt.Printf("Restored to: %s\n", record.Path)
	return f.reg.Save()
}
