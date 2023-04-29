package packagemanagers

import (
	"encoding/json"
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

// IsInstalled returns true if the package manager is installed.
func (scoopPackageManager *ScoopPackageManager) IsInstalled() (bool, error) {
	fmt.Printf("Checking if package manager \"%s\" is already installed...\n", scoopPackageManager.Name())

	_, err := shell.RunShellCommand(scoopPackageManager.Name(), false, nil, "--version")
	if err != nil {
		return false, nil
	}

	return true, nil
}

// Install installs the package manager.
func (scoopPackageManager *ScoopPackageManager) Install() error {
	fmt.Printf("Installing package manager \"%s\"...\n", scoopPackageManager.Name())

	_, err := shell.RunShellCommand("powershell", true, nil, "irm get.scoop.sh | iex")
	if err != nil {
		return err
	}

	return nil
}

// Uninstall uninstalls the package manager.
func (scoopPackageManager *ScoopPackageManager) Uninstall() error {
	fmt.Printf("Uninstalling package manager \"%s\"...\n", scoopPackageManager.Name())

	// Scoop doesn't return non-zero exit codes, so we have to check the output to see if the operation was successful.
	// Todo: Add a regex to make sure the operation was successful.

	_, err := shell.RunShellCommand(scoopPackageManager.Name(), true, nil, "uninstall", scoopPackageManager.Name())
	if err != nil {
		return err
	}

	return nil
}

// InstalledPackages returns a slice containing information about all packages that are installed.
func (scoopPackageManager *ScoopPackageManager) InstalledPackages() ([]*Package, error) {
	fmt.Printf("Getting installed package information from package manager \"%s\"...\n", scoopPackageManager.Name())

	outputJsonRegex, err := regexp.Compile("(?s)(.*)")
	if err != nil {
		return nil, err
	}

	outputJson, err := shell.RunShellCommand(scoopPackageManager.Name(), false, outputJsonRegex, "export")
	if err != nil {
		return nil, err
	}

	type ScoopExport struct {
		Apps []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"apps"`
	}
	var scoopExport ScoopExport

	err = json.Unmarshal([]byte(outputJson), &scoopExport)
	if err != nil {
		return nil, err
	}

	var installedPackages []*Package
	for _, app := range scoopExport.Apps {
		installedPackages = append(installedPackages, &Package{
			Name:    app.Name,
			Version: &Version{VersionString: app.Version},
		})
	}

	return installedPackages, nil
}

// InstallPackage installs the package of the given name. If a version is given, that specific version of the package is
// installed. Otherwise, the latest version is installed.
//
// It returns the version of the package that was installed.
func (scoopPackageManager *ScoopPackageManager) InstallPackage(packageName string, version *Version) (*Version, error) {
	fmt.Printf("Installing package \"%s\"...\n", packageName)

	regexString := fmt.Sprintf("'%s' \\((.*)\\) was installed", packageName)
	versionCaptureRegex, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}

	capturedVersion, err := shell.RunShellCommand(scoopPackageManager.Name(), true, versionCaptureRegex, "install",
		packageName)
	if err != nil || capturedVersion == "" {
		if err == nil {
			err = fmt.Errorf("error installing package")
		}
		return nil, err
	}

	return &Version{VersionString: capturedVersion}, nil
}

// UpdatePackage updates the package of the given name. If a version is given, that specific version of the package is
// installed. Otherwise, the latest version is installed.
//
// It returns the version of the package that was installed.
func (scoopPackageManager *ScoopPackageManager) UpdatePackage(packageName string, version *Version) (*Version, error) {
	fmt.Printf("Updating package manager \"%s\"...\n", scoopPackageManager.Name())
	if _, err := shell.RunShellCommand(scoopPackageManager.Name(), false, nil, "update"); err != nil {
		return nil, err
	}

	fmt.Printf("Updating package \"%s\"...\n", packageName)

	regexString := fmt.Sprintf("'%s' \\((.*)\\) was installed", packageName)
	versionCaptureRegex, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}

	capturedVersion, err := shell.RunShellCommand(scoopPackageManager.Name(), true, versionCaptureRegex, "update",
		packageName)
	if err != nil || capturedVersion == "" {
		if err == nil {
			err = fmt.Errorf("error updating package")
		}
		return nil, err
	}

	return &Version{VersionString: capturedVersion}, nil
}

// UninstallPackage uninstalls the package of the given name.
func (scoopPackageManager *ScoopPackageManager) UninstallPackage(packageName string) error {
	fmt.Printf("Uninstalling package \"%s\"...\n", packageName)

	// Scoop doesn't return non-zero exit codes, so we have to check the output to see if the operation was successful.
	regexString := fmt.Sprintf("('%s' was uninstalled)", packageName)
	successRegex, err := regexp.Compile(regexString)
	if err != nil {
		return err
	}

	capturedSuccess, err := shell.RunShellCommand(scoopPackageManager.Name(), true, successRegex, "uninstall",
		packageName)
	if err != nil || capturedSuccess == "" {
		if err == nil {
			err = fmt.Errorf("error uninstalling package")
		}
		return err
	}

	return nil
}
