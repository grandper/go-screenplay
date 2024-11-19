package question

// CreateFakeExistingFile creates a fake existing file.
func CreateFakeExistingFile(filePath string) *File {
	return newFile(filePath, true)
}

// CreateFakeMissingFile creates a fake missing file.
func CreateFakeMissingFile(filePath string) *File {
	return newFile(filePath, false)
}

// CreateFakeExistingDirectory creates a fake existing directory.
func CreateFakeExistingDirectory(directoryPath string) *Directory {
	return newDirectory(directoryPath, true)
}

// CreateFakeMissingDirectory creates a fake missing directory.
func CreateFakeMissingDirectory(directoryPath string) *Directory {
	return newDirectory(directoryPath, false)
}
