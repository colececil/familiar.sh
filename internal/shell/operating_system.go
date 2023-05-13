package shell

import "runtime"

// IsWindows returns whether the current operating system is Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
