package register

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

const (
	Filename    = "trash.registry.json"
	DefaultPath = ".Trash"
)

type Register struct {
	path    string
	records []Record
}

func New(path string) (*Register, error) {
	if path == "" {
		path = DefaultPath
	}

	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, path[1:])
	}

	if !filepath.IsAbs(path) {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, path)
	}

	return &Register{
		path:    filepath.Join(path, Filename),
		records: make([]Record, 0),
	}, nil
}

func (r *Register) Load() error {
	file, err := os.Open(r.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&r.records)
}

func (r *Register) Save() error {
	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(r.path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(r.records)
}

func (r *Register) NewID() int {
	if len(r.records) == 0 {
		return 1
	}
	last := r.records[len(r.records)-1]
	return last.ID + 1
}

func (r *Register) Add(name, path string) (Record, error) {
	record := Record{
		ID:   r.NewID(),
		Name: name,
		Path: path,
	}
	r.records = append(r.records, record)
	return record, nil
}

func (r *Register) Remove(id int) error {
	for i, record := range r.records {
		if record.ID == id {
			r.records = slices.Delete(r.records, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("record with ID %d not found", id)
}

func (r *Register) Get(id int) (Record, error) {
	for _, record := range r.records {
		if record.ID == id {
			return record, nil
		}
	}
	return Record{}, fmt.Errorf("record with ID %d not found", id)
}

func (r *Register) List() []Record {
	return r.records
}

func EncodePath(filePath string) (string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	baseName := filepath.Base(absPath)
	dirPath := filepath.Dir(absPath)

	encodedDir := base64.StdEncoding.EncodeToString([]byte(dirPath))

	return fmt.Sprintf("%s__%s#0", baseName, encodedDir), nil
}

func DecodePath(encodedFile string) (string, error) {
	encodedFile = strings.Split(encodedFile, "#")[0]

	parts := strings.SplitN(encodedFile, "__", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid encoded file format: %s", encodedFile)
	}

	baseName := parts[0]
	encodedPath := parts[1]

	dirPathBytes, err := base64.StdEncoding.DecodeString(encodedPath)
	if err != nil {
		return "", err
	}

	return filepath.Join(string(dirPathBytes), baseName), nil
}

func GetTrashRoot(filePath string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	isMounted := false
	for _, prefix := range []string{
		"/media/",
		"/run/media/",
		"/mnt/",
		"/Volumes/",
	} {
		if strings.HasPrefix(absPath, prefix) {
			isMounted = true
			break
		}
	}

	if isMounted {
		cmd := exec.Command("df", "-P", absPath)
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			if len(lines) >= 2 {
				fields := strings.Fields(lines[1])
				if len(fields) >= 6 {
					mountPoint := fields[5]
					return filepath.Join(mountPoint, DefaultPath), nil
				}
			}
		}
	}

	return filepath.Join(home, DefaultPath), nil
}
