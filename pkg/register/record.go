package register

type Record struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
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
