package test

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/system"
	"io"
)

// FileSystemServiceDouble is an implementation of system.FileSystemService to use as a test double.
type FileSystemServiceDouble struct {
	system.FileSystemService
	xdgConfigHome     string
	expectedFilePaths map[string]string
	createdFiles      []*FileDouble
}

// NewFileSystemServiceDouble returns a new instance of FileSystemServiceDouble.
func NewFileSystemServiceDouble() *FileSystemServiceDouble {
	return &FileSystemServiceDouble{
		expectedFilePaths: make(map[string]string),
	}
}

// SetXdgConfigHome sets the XDG config home directory.
func (double *FileSystemServiceDouble) SetXdgConfigHome(xdgConfigHome string) {
	double.xdgConfigHome = xdgConfigHome
}

// GetXdgConfigHome implements system.FileSystemService.GetXdgConfigHome by returning the value of the xdgConfigHome
// string.
func (double *FileSystemServiceDouble) GetXdgConfigHome() string {
	return double.xdgConfigHome
}

// SetFileContentForExpectedPath sets the file content to return when calling ReadFile with the given path.
func (double *FileSystemServiceDouble) SetFileContentForExpectedPath(path string, content string) {
	double.expectedFilePaths[path] = content
}

// ReadFile implements system.FileSystemService.ReadFile by returning the value in the expectedFilePaths map for
// the given path, or an error if the path is not present.
func (double *FileSystemServiceDouble) ReadFile(path string) ([]byte, error) {
	content, isPresent := double.expectedFilePaths[path]
	if !isPresent {
		return nil, fmt.Errorf("unexpected file path '%s'", path)
	}
	return []byte(content), nil
}

// CreateFile implements system.FileSystemService.CreateFile by creating a new FileDouble for the given path and adding
// it to the createdFiles slice.
func (double *FileSystemServiceDouble) CreateFile(path string) (io.WriteCloser, error) {
	file := NewFileDouble(path)
	double.createdFiles = append(double.createdFiles, file)
	return file, nil
}

// GetCreatedFiles returns a slice containing all file doubles that have been created by CreateFile.
func (double *FileSystemServiceDouble) GetCreatedFiles() []*FileDouble {
	return double.createdFiles
}
