package test

import "github.com/colececil/familiar.sh/internal/system"

// OperatingSystemServiceDouble contains a test double for system.OperatingSystemService. The actual test double can be
// accessed via its OperatingSystemService field.
type OperatingSystemServiceDouble struct {
	OperatingSystemService *system.OperatingSystemService
}

var isWindows bool

// NewOperatingSystemServiceDouble returns a new instance of OperatingSystemServiceDouble.
func NewOperatingSystemServiceDouble() *OperatingSystemServiceDouble {
	isWindows = false
	return &OperatingSystemServiceDouble{
		OperatingSystemService: system.NewOperatingSystemService(isWindowsFuncDouble),
	}
}

// SetIsWindows sets the value that will be returned by the test double's IsWindows function.
func (operatingSystemServiceDouble *OperatingSystemServiceDouble) SetIsWindows(newValue bool) {
	isWindows = newValue
}

// isWindowsFuncDouble is the implementation for the test double's IsWindows function.
func isWindowsFuncDouble() bool {
	return isWindows
}
