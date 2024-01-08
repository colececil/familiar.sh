package system

import (
	"errors"
	"io"
	"io/fs"
	"os"
)

// XdgConfigHomeGetter is an interface that can get the XDG config home directory.
type XdgConfigHomeGetter interface {
	GetXdgConfigHome() string
}

// XdgConfigHomeGetterFunc is a function that implements XdgConfigHomeGetter.
type XdgConfigHomeGetterFunc func() string

// GetXdgConfigHome implements XdgConfigHomeGetter.GetXdgConfigHome by calling the underlying function.
func (f XdgConfigHomeGetterFunc) GetXdgConfigHome() string {
	return f()
}

// AbsPathConverter is an interface that can convert a path to an absolute representation.
type AbsPathConverter interface {
	// Abs returns the absolute representation of the given path.
	Abs(path string) (string, error)
}

// AbsPathConverterFunc is a function that implements AbsPathConverter.
type AbsPathConverterFunc func(path string) (string, error)

// Abs implements AbsPathConverter.Abs by calling the underlying function.
func (f AbsPathConverterFunc) Abs(path string) (string, error) {
	return f(path)
}

// PathDirGetter is an interface that can get the directory of a path.
type PathDirGetter interface {
	Dir(path string) string
}

// PathDirGetterFunc is a function that implements PathDirGetter.
type PathDirGetterFunc func(path string) string

// Dir implements PathDirGetter.Dir by calling the underlying function.
func (f PathDirGetterFunc) Dir(path string) string {
	return f(path)
}

// FileExtensionGetter is an interface that can get the extension of a file specified by a path.
type FileExtensionGetter interface {
	// Ext returns the file extension of the given path. If the path has no extension, an empty string is returned.
	Ext(path string) string
}

// FileExtensionGetterFunc is a function that implements FileExtensionGetter.
type FileExtensionGetterFunc func(path string) string

// Ext implements FileExtensionGetter.Ext by calling the underlying function.
func (f FileExtensionGetterFunc) Ext(path string) string {
	return f(path)
}

// DirCreator is an interface that can create a directory at a given path.
type DirCreator interface {
	// CreateDirectory creates the directory at the given path, along with any required parent directories. All
	// directories created will have the given permissions. If the directory already exists, nothing happens.
	CreateDirectory(path string, permissions os.FileMode) error
}

// DirCreatorFunc is a function that implements DirCreator.
type DirCreatorFunc func(path string, permissions os.FileMode) error

// CreateDirectory implements DirCreator.CreateDirectory by calling the underlying function.
func (f DirCreatorFunc) CreateDirectory(path string, permissions os.FileMode) error {
	return f(path, permissions)
}

// FileExistenceChecker is an interface that can check the existence of a file specified by a path.
type FileExistenceChecker interface {
	// FileExists returns whether the file or directory at the given path exists.
	FileExists(path string) (bool, error)
}

// FileExistenceCheckerFunc is a function that implements FileExistenceChecker.
type FileExistenceCheckerFunc func(path string) (bool, error)

// FileExists implements FileExistenceChecker.FileExists by calling the underlying function.
func (f FileExistenceCheckerFunc) FileExists(path string) (bool, error) {
	return f(path)
}

// FileExists returns whether the file or directory at the given path exists.
func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// FileReader is an interface that can read a file specified by a path.
type FileReader interface {
	// ReadFile reads the file at the given path and returns the contents of the file.
	ReadFile(path string) ([]byte, error)
}

// FileReaderFunc is a function that implements FileReader.
type FileReaderFunc func(path string) ([]byte, error)

// ReadFile implements FileReader.ReadFile by calling the underlying function.
func (f FileReaderFunc) ReadFile(path string) ([]byte, error) {
	return f(path)
}

// FileCreator is an interface that can create a file at a given path.
type FileCreator interface {
	// CreateFile creates a file at the given path. If the file already exists, it is overwritten.
	CreateFile(path string) (io.WriteCloser, error)
}

// FileCreatorFunc is a function that implements FileCreator.
type FileCreatorFunc func(path string) (io.WriteCloser, error)

// CreateFile implements FileCreator.CreateFile by calling the underlying function.
func (f FileCreatorFunc) CreateFile(path string) (io.WriteCloser, error) {
	return f(path)
}
