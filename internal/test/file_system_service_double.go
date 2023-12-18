package test

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/system"
)

type FileSystemServiceDouble struct {
	system.FileSystemService
	xdgConfigHome     string
	expectedFilePaths map[string]string
}

func NewFileSystemServiceDouble() *FileSystemServiceDouble {
	return &FileSystemServiceDouble{
		expectedFilePaths: make(map[string]string),
	}
}

// SetXdgConfigHome sets the XDG config home directory.
func (double *FileSystemServiceDouble) SetXdgConfigHome(xdgConfigHome string) {
	double.xdgConfigHome = xdgConfigHome
}

// SetFileContentForExpectedPath sets the file content to return when calling ReadFile with the given path.
func (double *FileSystemServiceDouble) SetFileContentForExpectedPath(path string, content string) {
	double.expectedFilePaths[path] = content
}

func (double *FileSystemServiceDouble) GetXdgConfigHome() string {
	return double.xdgConfigHome
}

func (double *FileSystemServiceDouble) ReadFile(path string) ([]byte, error) {
	content, isPresent := double.expectedFilePaths[path]
	if !isPresent {
		return nil, fmt.Errorf("unexpected file path '%s'", path)
	}
	return []byte(content), nil
}
