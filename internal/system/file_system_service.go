package system

import (
	"github.com/adrg/xdg"
	"os"
)

type FileSystemService interface {
	GetXdgConfigHome() string
	CreateDirectory(path string, permissions os.FileMode) error
	FileExists(path string) (bool, error)
	ReadFile(path string) ([]byte, error)
	CreateFile(path string) (*os.File, error)
}

type ConcreteFileSystemService struct {
}

// NewConcreteFileSystemService creates a new instance of ConcreteFileSystemService.
func NewConcreteFileSystemService() *ConcreteFileSystemService {
	return &ConcreteFileSystemService{}
}

// GetXdgConfigHome returns the config home directory, according to the XDG specification.
func (fileSystemService *ConcreteFileSystemService) GetXdgConfigHome() string {
	return xdg.ConfigHome
}

// CreateDirectory creates the directory at the given path, along with any required parent directories. All directories
// created will have the given permissions. If the directory already exists, nothing happens.
func (fileSystemService *ConcreteFileSystemService) CreateDirectory(path string, permissions os.FileMode) error {
	return os.MkdirAll(path, permissions)
}

// FileExists returns whether the file at the given path exists.
func (fileSystemService *ConcreteFileSystemService) FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ReadFile reads the file at the given path.
func (fileSystemService *ConcreteFileSystemService) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// CreateFile creates a file at the given path. If the file already exists, it is overwritten.
func (fileSystemService *ConcreteFileSystemService) CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}
