package register

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

type DirSize struct {
	FolderName string
	Size       int64
	Seen       bool
	MTime      int64
}

const (
	Filename         = "trash.registry.json"
	DefaultPath      = "Trash/files"
	DefaultInfoPath  = "Trash/info"
	DefaultTrashPath = "Trash"
)

func GetUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Error: $HOME env variable is not set")
	}
	return home
}

func GetDataHome() string {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		return xdgDataHome
	}

	return GetUserHomeDir() + "/.local/share"
}

type Register struct {
	path    string
	records []Record
}

func New(path string) (*Register, error) {
	if path == "" {
		path = filepath.Join(GetDataHome(), DefaultTrashPath)
	}

	if strings.HasPrefix(path, "~") {
		path = filepath.Join(GetUserHomeDir(), path[1:])
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(GetUserHomeDir(), path)
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

func (r *Register) Add(name, encodedPath, path string) (Record, error) {
	record := Record{
		ID:          r.NewID(),
		Name:        name,
		Path:        path,
		EncodedPath: encodedPath,
		Info: RecordInfo{
			Path:      path,
			DeletedAt: time.Now().Format("2006-01-02"),
		},
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

	ts, err := getFileCreationTime(filePath)
	if err != nil {
		ts = time.Now()
	}

	return fmt.Sprintf("%s__%s__%d#0", baseName, encodedDir, ts.UnixNano()), nil
}

func (r *Register) GetInfoContent(ID int) (string, error) {
	record, err := r.Get(ID)
	if err != nil {
		return "", err
	}

	t, err := time.Parse("2006-01-02", record.Info.DeletedAt)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n", record.Path, t.Format("20040831T15:04:05")), nil
}

func (r *Register) GetDirSizeContent(path string, size int64) string {
	return fmt.Sprintf("%d %s\n", size, filepath.Base(path))
}

func (r *Register) GetDirSizes() ([]DirSize, error) {
	dirSizePath := filepath.Join(GetDirSizeRoot(), "directorysizes")
	if _, err := os.Stat(dirSizePath); os.IsNotExist(err) {
		return []DirSize{}, nil
	}

	file, err := os.Open(dirSizePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dirSizes []DirSize

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 3)
		if len(parts) != 3 {
			continue
		}

		size, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue
		}

		mtime, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue
		}

		name := parts[2]
		dirSizes = append(dirSizes, DirSize{
			FolderName: name,
			Size:       size,
			MTime:      mtime,
			Seen:       true,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dirSizes, nil
}

func GetTrashRoot() string {
	return filepath.Join(GetDataHome(), DefaultPath)
}

func GetTrashInfoRoot() string {
	return filepath.Join(GetDataHome(), DefaultInfoPath)
}

func GetDirSizeRoot() string {
	return filepath.Join(GetDataHome(), DefaultTrashPath)
}

func getFileCreationTime(filePath string) (time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}
