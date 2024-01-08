package test

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// FileSystemDouble is an implementation of the file system interfaces defined in the system package, to use as a test
// double. Paths are assumed to be in the Linux format.
type FileSystemDouble struct {
	xdgConfigHome       string
	createdFiles        map[string]*FileDouble
	errorMethodsAndArgs map[string]string
}

// NewFileSystemDouble returns a new instance of FileSystemDouble.
func NewFileSystemDouble() *FileSystemDouble {
	return &FileSystemDouble{
		createdFiles:        make(map[string]*FileDouble),
		errorMethodsAndArgs: make(map[string]string),
	}
}

// SetXdgConfigHome sets the XDG config home directory.
func (d *FileSystemDouble) SetXdgConfigHome(xdgConfigHome string) {
	d.xdgConfigHome = xdgConfigHome
}

// GetXdgConfigHome implements system.XdgConfigHomeGetter by returning the value of the xdgConfigHome string.
func (d *FileSystemDouble) GetXdgConfigHome() string {
	return d.xdgConfigHome
}

// Abs implements system.AbsPathConverter by returning the given path if it is absolute, or an error if it is not.
func (d *FileSystemDouble) Abs(path string) (string, error) {
	if arg, _ := d.errorMethodsAndArgs["Abs"]; arg == path {
		return "", errors.New("error getting absolute path")
	}
	if !strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("path \"%s\" is not absolute", path)
	}
	return path, nil
}

// Dir implements system.PathDirGetter by returning the directory of the given path. It assumes the path is absolute.
func (d *FileSystemDouble) Dir(path string) string {
	pathParts := strings.Split(path, "/")
	if len(pathParts) == 1 {
		return "/"
	}
	return strings.Join(pathParts[:len(pathParts)-1], "/")
}

// Ext implements system.FileExtensionGetter by returning the extension of the final part of the given path.
func (d *FileSystemDouble) Ext(path string) string {
	pathParts := strings.Split(path, "/")
	finalPathPart := pathParts[len(pathParts)-1]
	dotIndex := strings.LastIndex(finalPathPart, ".")
	if dotIndex == -1 {
		return ""
	}
	return finalPathPart[dotIndex:]
}

func (d *FileSystemDouble) CreateDirectory(path string, permissions os.FileMode) error {
	if arg, _ := d.errorMethodsAndArgs["CreateDirectory"]; arg == path {
		return errors.New("error creating directory")
	}
	exists, err := d.FileExists(path)
	if err != nil {
		return err
	}
	if !exists {
		file, err := d.CreateFile(path)
		if err != nil {
			return err
		}
		_ = file.Close()
	}
	return nil
}

// FileExists implements system.FileExistenceChecker by checking if the path exists in the createdFiles map.
func (d *FileSystemDouble) FileExists(path string) (bool, error) {
	if arg, _ := d.errorMethodsAndArgs["FileExists"]; arg == path {
		return false, errors.New("error checking if file exists")
	}
	_, isPresent := d.createdFiles[path]
	return isPresent, nil
}

// ReadFile implements system.FileReader by returning the content of the file at the given path. If the file does not
// exist, an error is returned.
func (d *FileSystemDouble) ReadFile(path string) ([]byte, error) {
	if arg, _ := d.errorMethodsAndArgs["ReadFile"]; arg == path {
		return nil, errors.New("error reading file")
	}
	file, isPresent := d.createdFiles[path]
	if !isPresent {
		return nil, fs.ErrNotExist
	}
	if !file.isClosed {
		return nil, fmt.Errorf("file at path \"%s\" is not closed", path)
	}
	return []byte(file.content), nil
}

// CreateFile implements system.FileCreator by creating a new FileDouble for the given path and adding it to the
// createdFiles map.
func (d *FileSystemDouble) CreateFile(path string) (io.WriteCloser, error) {
	if arg, _ := d.errorMethodsAndArgs["CreateFile"]; arg == path {
		return nil, errors.New("error creating file")
	}
	file := NewFileDouble(path)
	d.createdFiles[path] = file
	return file, nil
}

// GetCreatedFile returns a reference to the FileDouble created for the given path, if it exists. It also returns a
// boolean indicating whether the FileDouble exists.
func (d *FileSystemDouble) GetCreatedFile(path string) (*FileDouble, bool) {
	file, isPresent := d.createdFiles[path]
	return file, isPresent
}

// ReturnErrorFromMethod tells the FileSystemDouble to return an error any time the given method is called with
// the given argument. Valid method names are "Abs", "CreateDirectory", "FileExists", "ReadFile", and "CreateFile". If
// the given method does not take an argument, the argument can be nil or an empty string.
func (d *FileSystemDouble) ReturnErrorFromMethod(methodName string, arg string) {
	d.errorMethodsAndArgs[methodName] = arg
}
