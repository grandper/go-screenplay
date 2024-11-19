package question

import "path"

// Directory represents a directory.
type Directory struct {
	dir    string
	name   string
	exists bool
}

func newDirectory(directoryPath string, exists bool) *Directory {
	return &Directory{
		dir:    path.Dir(directoryPath),
		name:   path.Base(directoryPath),
		exists: exists,
	}
}

// Dir returns the dir to the directory.
func (d *Directory) Dir() string {
	return d.dir
}

// Name returns the name of the file.
func (d *Directory) Name() string {
	return d.name
}

// Exists returns true if the file exists.
func (d *Directory) Exists() bool {
	return d.exists
}
