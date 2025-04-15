package register

type RecordInfo struct {
	Path      string `json:"path"`
	DeletedAt string `json:"deleted_at"`
}

type Record struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	EncodedPath string     `json:"encoded_path"`
	Info        RecordInfo `json:"info"`
}

func (r *Record) String() string {
	return r.Name
}

func (r *Record) IsDir() bool {
	return r.Path[len(r.Path)-1] == '/'
}

func (r *Record) IsFile() bool {
	return !r.IsDir()
}
