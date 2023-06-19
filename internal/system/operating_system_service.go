package system

import "runtime"

// OperatingSystemService provides information about the operating system.
type OperatingSystemService struct {
}

// NewOperatingSystemService returns a new instance of OperatingSystemService.
func NewOperatingSystemService() *OperatingSystemService {
	return &OperatingSystemService{}
}

// IsWindows returns whether the current operating system is Windows.
func (operatingSystemService *OperatingSystemService) IsWindows() bool {
	return runtime.GOOS == "windows"
}
