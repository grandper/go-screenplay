package fixture

import "github.com/grandper/go-screenplay/v2/extensions/filesystem/question"

// CreateFakeExistingFile creates a fake existing file.
func CreateFakeExistingFile(filePath string) *question.File {
	return question.NewFile(filePath, true)
}

// CreateFakeMissingFile creates a fake missing file.
func CreateFakeMissingFile(filePath string) *question.File {
	return question.NewFile(filePath, false)
}

// CreateFakeExistingDirectory creates a fake existing directory.
func CreateFakeExistingDirectory(directoryPath string) *question.Directory {
	return question.NewDirectory(directoryPath, true)
}

// CreateFakeMissingDirectory creates a fake missing directory.
func CreateFakeMissingDirectory(directoryPath string) *question.Directory {
	return question.NewDirectory(directoryPath, false)
}
