package test

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/system"
	"io"
	"os"
	"strings"
)

// FileSystemServiceDouble is an implementation of system.FileSystemService to use as a test double. Paths are assumed
// to be in the Linux format.
type FileSystemServiceDouble struct {
	system.FileSystemService
	xdgConfigHome string
	createdFiles  map[string]*FileDouble
}

// NewFileSystemServiceDouble returns a new instance of FileSystemServiceDouble.
func NewFileSystemServiceDouble() *FileSystemServiceDouble {
	return &FileSystemServiceDouble{
		createdFiles: make(map[string]*FileDouble),
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

// Abs implements system.FileSystemService.Abs by returning the given path if it is absolute, or an error if it is
// not.
func (double *FileSystemServiceDouble) Abs(path string) (string, error) {
	if !strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("path \"%s\" is not absolute", path)
	}
	return path, nil
}

// Dir implements system.FileSystemService.Dir by returning the directory of the given path. It assumes the path is
// absolute.
func (double *FileSystemServiceDouble) Dir(path string) string {
	pathParts := strings.Split(path, "/")
	if len(pathParts) == 1 {
		return "/"
	}
	return strings.Join(pathParts[:len(pathParts)-1], "/")
}

// Ext implements system.FileSystemService.Ext by returning the extension of the final part of the given path.
func (double *FileSystemServiceDouble) Ext(path string) string {
	pathParts := strings.Split(path, "/")
	finalPathPart := pathParts[len(pathParts)-1]
	dotIndex := strings.LastIndex(finalPathPart, ".")
	if dotIndex == -1 {
		return ""
	}
	return finalPathPart[dotIndex:]
}

func (double *FileSystemServiceDouble) CreateDirectory(path string, permissions os.FileMode) error {
	exists, err := double.FileExists(path)
	if err != nil {
		return err
	}
	if !exists {
		file, err := double.CreateFile(path)
		defer func(closer io.Closer) {
			_ = closer.Close()
		}(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// FileExists implements system.FileSystemService.FileExists by checking if the path exists in the createdFiles
// map.
func (double *FileSystemServiceDouble) FileExists(path string) (bool, error) {
	_, isPresent := double.createdFiles[path]
	return isPresent, nil
}

// ReadFile implements system.FileSystemService.ReadFile by returning the content of the file at the given path. If the
// file does not exist, an error is returned.
func (double *FileSystemServiceDouble) ReadFile(path string) ([]byte, error) {
	file, isPresent := double.createdFiles[path]
	if !isPresent {
		return nil, fmt.Errorf("file at path \"%s\" does not exist", path)
	}
	if !file.isClosed {
		return nil, fmt.Errorf("file at path \"%s\" is not closed", path)
	}
	return []byte(file.content), nil
}

// CreateFile implements system.FileSystemService.CreateFile by creating a new FileDouble for the given path and adding
// it to the createdFiles map.
func (double *FileSystemServiceDouble) CreateFile(path string) (io.WriteCloser, error) {
	file := NewFileDouble(path)
	double.createdFiles[path] = file
	return file, nil
}

// GetCreatedFile returns a reference to the FileDouble created for the given path, if it exists. It also returns a
// boolean indicating whether the FileDouble exists.
func (double *FileSystemServiceDouble) GetCreatedFile(path string) (*FileDouble, bool) {
	file, isPresent := double.createdFiles[path]
	return file, isPresent
}
