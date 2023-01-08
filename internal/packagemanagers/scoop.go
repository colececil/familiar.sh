package packagemanagers

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/shell"
	"regexp"
)

type ScoopPackageManager struct {
}

// Name returns the name of the package manager.
func (scoopPackageManager *ScoopPackageManager) Name() string {
	return "scoop"
}

// InstallPackage installs the package of the given name. If a version is given, that specific version of the package is
// installed. Otherwise, the latest version is installed.
//
// It returns the version of the package that was installed.
func (scoopPackageManager *ScoopPackageManager) InstallPackage(packageName string, version *string) (string, error) {
	regexString := fmt.Sprintf("'%s' \\((.*)\\) was installed", packageName)
	versionCaptureRegex, err := regexp.Compile(regexString)
	if err != nil {
		return "", err
	}

	capturedVersion, err := shell.RunShellCommand("scoop", versionCaptureRegex, "install", packageName)
	if err != nil {
		return "", err
	}

	return capturedVersion, nil
}

// UninstallPackage uninstalls the package of the given name.
func (scoopPackageManager *ScoopPackageManager) UninstallPackage(packageName string) error {
	_, err := shell.RunShellCommand("scoop", nil, "uninstall", packageName)
	return err
}
