package question

import "path"

// File represents a file.
type File struct {
	dir    string
	name   string
	exists bool
}

func newFile(filePath string, exists bool) *File {
	return &File{
		dir:    path.Dir(filePath),
		name:   path.Base(filePath),
		exists: exists,
	}
}

// Dir returns the dir to the file.
func (f *File) Dir() string {
	return f.dir
}

// Name returns the name of the file.
func (f *File) Name() string {
	return f.name
}

// Exists returns true if the file exists.
func (f *File) Exists() bool {
	return f.exists
}
