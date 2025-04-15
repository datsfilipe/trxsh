package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/datsfilipe/trxsh/pkg/register"
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

		trashRoot := register.GetTrashRoot()
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

		record, err := c.reg.Add(filepath.Base(file), absPath)
		if err != nil {
			return err
		}

		err = c.SaveTrashInfo(record.ID, encodedName)
		if err != nil {
			return err
		}

		if err := os.Rename(file, trashPath); err != nil {
			return err
		}

		if err := c.CalcDirSize(trashRoot); err != nil {
			return err
		}
	}

	return c.reg.Save()
}

func (c *CLI) RawList() []register.Record {
	return c.reg.List()
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

	trashRoot := register.GetTrashRoot()
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
	if err := c.DeleteTrashInfo(record.ID, encodedName); err != nil {
		return err
	}

	if err := c.reg.Remove(record.ID); err != nil {
		return err
	}

	if err := c.CalcDirSize(trashRoot); err != nil {
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

func (c *CLI) SaveTrashInfo(ID int, encodedName string) error {
	infoRoot := register.GetTrashInfoRoot()
	if err := os.MkdirAll(infoRoot, 0755); err != nil {
		return err
	}
	content, err := c.reg.GetInfoContent(ID)
	if err != nil {
		return err
	}

	infoPath := filepath.Join(infoRoot, encodedName+".trashinfo")
	return os.WriteFile(infoPath, []byte(content), 0644)
}

func (c *CLI) DeleteTrashInfo(ID int, encodedName string) error {
	infoRoot := register.GetTrashInfoRoot()
	infoPath := filepath.Join(infoRoot, encodedName+".trashinfo")

	return os.Remove(infoPath)
}

func (c *CLI) PrintDirSizes() error {
	dirSizes, err := c.reg.GetDirSizes()
	if err != nil {
		return err
	}

	for _, dirSize := range dirSizes {
		sizeInMB := float64(dirSize.Size) / 1024 / 1024
		fmt.Printf("%s (%.2f MB)\n", dirSize.FolderName, sizeInMB)
	}

	return nil
}

func (c *CLI) CalcDirSize(path string) error {
	infoDir := register.GetTrashInfoRoot()
	existingDirSizes, err := c.reg.GetDirSizes()
	if err != nil {
		existingDirSizes = []register.DirSize{}
	}

	existingDirSizesMap := make(map[string]register.DirSize)
	for _, dirSize := range existingDirSizes {
		existingDirSizesMap[dirSize.FolderName] = dirSize
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	dirSizes := make(map[string]register.DirSize)
	totalSize := int64(0)

	for _, file := range files {
		itemPath := filepath.Join(path, file.Name())
		fileInfo, err := os.Stat(itemPath)
		if err != nil {
			continue
		}

		if fileInfo.IsDir() {
			encodedName := file.Name()
			parts := strings.Split(encodedName, "__")
			baseName := parts[0]
			trashInfoPath := filepath.Join(infoDir, encodedName+".trashinfo")
			trashInfo, err := os.Stat(trashInfoPath)
			if err != nil {
				trashInfoPath = filepath.Join(infoDir, baseName+".trashinfo")
				trashInfo, err = os.Stat(trashInfoPath)
				if err != nil {
					continue
				}
			}

			content, err := os.ReadFile(trashInfoPath)
			if err != nil {
				continue
			}

			var folderName string
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "Path=") {
					retrievedPath := strings.TrimPrefix(line, "Path=")
					retrievedPath = strings.TrimSpace(retrievedPath)
					folderName = filepath.Base(retrievedPath)
					break
				}
			}

			trashInfoMtime := trashInfo.ModTime().Unix()
			existingDirSize, exists := existingDirSizesMap[encodedName]

			if !exists || existingDirSize.MTime != trashInfoMtime {
				size, err := c.CalculateDirSizeRecursively(itemPath)
				if err != nil {
					continue
				}
				totalSize += size
				dirSizes[encodedName] = register.DirSize{
					FolderName: folderName,
					Size:       size,
					MTime:      trashInfoMtime,
					Seen:       true,
				}
			} else {
				totalSize += existingDirSize.Size
				dirSizes[encodedName] = register.DirSize{
					FolderName: existingDirSize.FolderName,
					Size:       existingDirSize.Size,
					MTime:      existingDirSize.MTime,
					Seen:       true,
				}
			}
		} else {
			totalSize += fileInfo.Size()
		}
	}

	var dirSizesToSave []register.DirSize
	for _, dirSize := range dirSizes {
		dirSizesToSave = append(dirSizesToSave, dirSize)
	}

	if len(dirSizesToSave) == 0 {
		c.DeleteDirSize()
		return nil
	}

	return c.SaveDirSize(dirSizesToSave)
}

func (c *CLI) CalculateDirSizeRecursively(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

func (c *CLI) SaveDirSize(dirSizes []register.DirSize) error {
	dirSizeRoot := register.GetDirSizeRoot()
	if err := os.MkdirAll(dirSizeRoot, 0755); err != nil {
		return err
	}

	dirSizePath := filepath.Join(dirSizeRoot, "directorysizes")
	var content strings.Builder
	for _, dirSize := range dirSizes {
		line := fmt.Sprintf("%d %d %s\n", dirSize.MTime, dirSize.Size, dirSize.FolderName)
		content.WriteString(line)
	}

	return os.WriteFile(dirSizePath, []byte(content.String()), 0644)
}

func (c *CLI) DeleteDirSize() {
	dirSizeRoot := register.GetDirSizeRoot()
	dirSizePath := filepath.Join(dirSizeRoot, "directorysizes")

	if _, err := os.Stat(dirSizePath); err == nil {
		os.Remove(dirSizePath)
	}
}
