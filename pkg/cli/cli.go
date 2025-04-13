package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/datsfilipe/trash-cli/pkg/register"
)

type CLI struct {
	reg *register.Register
}

func New() (*CLI, error) {
	reg, err := register.New("")
	if err != nil {
		return nil, err
	}

	if err := reg.Load(); err != nil {
		return nil, err
	}

	return &CLI{reg: reg}, nil
}

func (c *CLI) Trash(args []string) error {
	if len(args) == 0 {
		return errors.New("no files specified")
	}

	for _, file := range args {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("File not found: %s\n", file)
			continue
		}

		trashRoot, err := register.GetTrashRoot(file)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(trashRoot, 0755); err != nil {
			return err
		}

		encodedName, err := register.EncodePath(file)
		if err != nil {
			return err
		}

		trashPath := filepath.Join(trashRoot, encodedName)

		absPath, err := filepath.Abs(file)
		if err != nil {
			return err
		}

		_, err = c.reg.Add(filepath.Base(file), absPath)
		if err != nil {
			return err
		}

		if err := os.Rename(file, trashPath); err != nil {
			return err
		}
	}

	return c.reg.Save()
}

func (c *CLI) List() error {
	records := c.reg.List()

	if len(records) == 0 {
		fmt.Println("No files in trash")
		return nil
	}

	for _, record := range records {
		fmt.Printf("%d: %s -> %s\n", record.ID, record.Name, record.Path)
	}

	return nil
}

func (c *CLI) Restore(idStr string) error {
	if idStr == "" {
		return errors.New("please provide an ID")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	record, err := c.reg.Get(id)
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

	if err := c.reg.Remove(record.ID); err != nil {
		return err
	}

	fmt.Printf("Restored to: %s\n", record.Path)

	return c.reg.Save()
}

func (c *CLI) Cleanup() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	trashDir := filepath.Join(home, ".Trash")
	if stat, err := os.Stat(trashDir); err == nil && stat.IsDir() {
		if err := os.RemoveAll(trashDir); err != nil {
			return fmt.Errorf("failed to remove trash directory %q: %w", trashDir, err)
		}
	}

	newReg, err := register.New("")
	if err != nil {
		return err
	}

	c.reg = newReg
	return c.reg.Save()
}
