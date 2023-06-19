package system

import "runtime"

// OperatingSystemService provides information about the operating system.
type OperatingSystemService struct {
	isWindowsFunc IsWindowsFunc
}

// NewOperatingSystemService returns a new instance of OperatingSystemService.
func NewOperatingSystemService(isWindows IsWindowsFunc) *OperatingSystemService {
	return &OperatingSystemService{
		isWindowsFunc: isWindows,
	}
}

// IsWindowsFunc is a function for determining whether the current operating system is Windows.
type IsWindowsFunc func() bool

// NewIsWindowsFunc returns a new function for determining whether the current operating system is Windows.
func NewIsWindowsFunc() IsWindowsFunc {
	return defaultIsWindowsFunc
}

// IsWindows returns whether the current operating system is Windows.
func (operatingSystemService *OperatingSystemService) IsWindows() bool {
	return operatingSystemService.isWindowsFunc()
}

// defaultIsWindowsFunc returns the default implementation of IsWindowsFunc.
func defaultIsWindowsFunc() bool {
	return runtime.GOOS == "windows"
}
