package system

import (
	"errors"
	"github.com/adrg/xdg"
	"io"
	"io/fs"
	"os"
)

// FileSystemService is an interface for a service that provides file system operations.
type FileSystemService interface {
	// GetXdgConfigHome returns the config home directory, according to the XDG specification.
	GetXdgConfigHome() string

	// CreateDirectory creates the directory at the given path, along with any required parent directories. All directories
	// created will have the given permissions. If the directory already exists, nothing happens.
	CreateDirectory(path string, permissions os.FileMode) error

	// FileExists returns whether the file at the given path exists.
	FileExists(path string) (bool, error)

	// ReadFile reads the file at the given path.
	ReadFile(path string) ([]byte, error)

	// CreateFile creates a file at the given path. If the file already exists, it is overwritten.
	CreateFile(path string) (io.WriteCloser, error)
}

// NewFileSystemService creates a new instance of FileSystemService.
func NewFileSystemService() FileSystemService {
	return &fileSystemService{}
}

type fileSystemService struct {
}

// GetXdgConfigHome is a concrete implementation of FileSystemService.GetXdgConfigHome.
func (service *fileSystemService) GetXdgConfigHome() string {
	return xdg.ConfigHome
}

// CreateDirectory is a concrete implementation of FileSystemService.CreateDirectory.
func (service *fileSystemService) CreateDirectory(path string, permissions os.FileMode) error {
	return os.MkdirAll(path, permissions)
}

// FileExists is a concrete implementation of FileSystemService.FileExists.
func (service *fileSystemService) FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ReadFile is a concrete implementation of FileSystemService.ReadFile.
func (service *fileSystemService) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// CreateFile is a concrete implementation of FileSystemService.CreateFile.
func (service *fileSystemService) CreateFile(path string) (io.WriteCloser, error) {
	return os.Create(path)
}
