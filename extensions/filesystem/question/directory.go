package question

import "path"

// Directory represents a directory.
type Directory struct {
	dir    string
	name   string
	exists bool
}

// NewDirectory creates the answer of DirectoryNamed: the directory at the path,
// which exists or not. Dir and Name split the path as it was written; the path
// is not resolved.
func NewDirectory(directoryPath string, exists bool) *Directory {
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

// Name returns the name of the directory.
func (d *Directory) Name() string {
	return d.name
}

// Exists returns true if the directory exists.
func (d *Directory) Exists() bool {
	return d.exists
}
