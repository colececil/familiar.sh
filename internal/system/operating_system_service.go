package system

// OperatingSystemService provides information about the operating system.
type OperatingSystemService struct {
	currentOperatingSystem OperatingSystem
}

type OperatingSystem string

const WindowsOperatingSystem OperatingSystem = "windows"
const MacOsOperatingSystem OperatingSystem = "darwin"
const LinuxOperatingSystem OperatingSystem = "linux"

// NewOperatingSystemService returns a new instance of OperatingSystemService.
func NewOperatingSystemService(currentOperatingSystem OperatingSystem) *OperatingSystemService {
	return &OperatingSystemService{
		currentOperatingSystem: currentOperatingSystem,
	}
}

// IsWindows returns whether the current operating system is Windows.
func (operatingSystemService *OperatingSystemService) IsWindows() bool {
	return operatingSystemService.currentOperatingSystem == WindowsOperatingSystem
}
