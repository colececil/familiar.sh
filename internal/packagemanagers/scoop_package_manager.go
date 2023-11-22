package packagemanagers

import (
	"encoding/json"
	"fmt"
	"github.com/colececil/familiar.sh/internal/system"
	"io"
	"regexp"
	"slices"
	"strings"
)

// ScoopPackageManager implements the PackageManager interface for the Scoop package manager.
type ScoopPackageManager struct {
	operatingSystemService *system.OperatingSystemService
	shellCommandService    *system.ShellCommandService
	outputWriter           io.Writer
}

// NewScoopPackageManager returns a new instance of ScoopPackageManager.
func NewScoopPackageManager(operatingSystemService *system.OperatingSystemService,
	shellCommandService *system.ShellCommandService, outputWriter io.Writer) *ScoopPackageManager {
	return &ScoopPackageManager{
		operatingSystemService: operatingSystemService,
		shellCommandService:    shellCommandService,
		outputWriter:           outputWriter,
	}
}

// Name returns the name of the package manager.
func (scoopPackageManager *ScoopPackageManager) Name() string {
	return "scoop"
}

// IsSupported returns whether the package manager is supported on the current machine.
func (scoopPackageManager *ScoopPackageManager) IsSupported() bool {
	return scoopPackageManager.operatingSystemService.IsWindows()
}

// IsInstalled returns true if the package manager is installed.
func (scoopPackageManager *ScoopPackageManager) IsInstalled() (bool, error) {
	_, err := fmt.Fprintf(scoopPackageManager.outputWriter, "Checking if package manager \"%s\" is installed...\n",
		scoopPackageManager.Name())
	if err != nil {
		return false, err
	}

	_, err = scoopPackageManager.shellCommandService.RunShellCommand(scoopPackageManager.Name(), false, nil,
		"--version")
	if err != nil {
		return false, nil
	}

	return true, nil
}

// Install installs the package manager.
func (scoopPackageManager *ScoopPackageManager) Install() error {
	_, err := fmt.Fprintf(scoopPackageManager.outputWriter, "Installing package manager \"%s\"...\n",
		scoopPackageManager.Name())
	if err != nil {
		return err
	}

	_, err = scoopPackageManager.shellCommandService.RunShellCommand("powershell", true, nil, "irm get.scoop.sh | iex")
	if err != nil {
		return err
	}

	return nil
}

// Update updates the package manager.
func (scoopPackageManager *ScoopPackageManager) Update() error {
	_, err := fmt.Fprintf(scoopPackageManager.outputWriter, "Updating package manager \"%s\"...\n",
		scoopPackageManager.Name())
	if err != nil {
		return err
	}

	// Scoop doesn't return non-zero exit codes, so we have to check the output to see if the operation was successful.
	regexString := "(Scoop was updated)"
	successRegex, err := regexp.Compile(regexString)
	if err != nil {
		return err
	}

	capturedSuccess, err := scoopPackageManager.shellCommandService.RunShellCommand(scoopPackageManager.Name(), true,
		successRegex, "update")
	if err != nil || capturedSuccess == "" {
		if err == nil {
			err = fmt.Errorf("error updating package manager \"%s\"", scoopPackageManager.Name())
		}
		return err
	}

	return nil
}

// Uninstall uninstalls the package manager.
func (scoopPackageManager *ScoopPackageManager) Uninstall() error {
	_, err := fmt.Fprintf(scoopPackageManager.outputWriter, "Uninstalling package manager \"%s\"...\n",
		scoopPackageManager.Name())
	if err != nil {
		return err
	}

	// Scoop doesn't return non-zero exit codes, so we have to check the output to see if the operation was successful.
	regexString := fmt.Sprintf("('%s' was uninstalled)", scoopPackageManager.Name())
	successRegex, err := regexp.Compile(regexString)
	if err != nil {
		return err
	}

	capturedSuccess, err := scoopPackageManager.shellCommandService.RunShellCommand(scoopPackageManager.Name(), true,
		successRegex, "uninstall", scoopPackageManager.Name())
	if err != nil || capturedSuccess == "" {
		if err == nil {
			err = fmt.Errorf("error uninstalling package manager \"%s\"", scoopPackageManager.Name())
		}
		return err
	}

	return nil
}

// InstalledPackages returns a slice containing information about all packages that are installed.
func (scoopPackageManager *ScoopPackageManager) InstalledPackages() ([]*Package, error) {
	_, err := fmt.Fprintf(scoopPackageManager.outputWriter,
		"Getting installed package information from package manager \"%s\"...\n", scoopPackageManager.Name())
	if err != nil {
		return nil, err
	}

	jsonCaptureRegex, err := regexp.Compile("(?s)(.*)")
	if err != nil {
		return nil, err
	}

	capturedJson, err := scoopPackageManager.shellCommandService.RunShellCommand(scoopPackageManager.Name(), false,
		jsonCaptureRegex, "export")
	if err != nil {
		return nil, err
	}

	type scoopExportData struct {
		Apps []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"apps"`
	}
	var scoopExport scoopExportData

	err = json.Unmarshal([]byte(capturedJson), &scoopExport)
	if err != nil {
		return nil, err
	}

	var installedPackages = make(map[string]*Package)
	for _, app := range scoopExport.Apps {
		installedPackages[app.Name] = NewPackage(app.Name, NewVersion(app.Version), NewVersion(app.Version))
	}

	packagesCaptureRegex, err := regexp.Compile("(?s)^.*----\\n(([^\\n]*(\\n)??)*)\\n*$")
	if err != nil {
		return nil, err
	}

	capturedPackages, err := scoopPackageManager.shellCommandService.RunShellCommand(scoopPackageManager.Name(), false,
		packagesCaptureRegex, "status")
	if err != nil {
		return nil, err
	}

	if capturedPackages != "" {
		for _, packageLine := range strings.Split(capturedPackages, "\n") {
			packageFields := strings.Fields(packageLine)

			if len(packageFields) != 3 {
				return nil, fmt.Errorf("unexpected number of fields in line: %s", packageLine)
			}

			installedPackage, isPresent := installedPackages[packageFields[0]]
			if isPresent {
				installedPackage.LatestVersion = NewVersion(packageFields[2])
			}
		}
	}

	installedPackagesSlice := make([]*Package, 0)
	for _, installedPackage := range installedPackages {
		installedPackagesSlice = append(installedPackagesSlice, installedPackage)
	}
	slices.SortFunc(installedPackagesSlice, func(package1, package2 *Package) int {
		return strings.Compare(
			strings.ToLower(package1.Name),
			strings.ToLower(package2.Name),
		)
	})

	return installedPackagesSlice, nil
}

// InstallPackage installs the package of the given name. If a version is given, that specific version of the package is
// installed. Otherwise, the latest version is installed.
//
// It returns the version of the package that was installed.
func (scoopPackageManager *ScoopPackageManager) InstallPackage(packageName string, version *Version) (*Version, error) {
	_, err := fmt.Fprintf(scoopPackageManager.outputWriter, "Installing package \"%s\"...\n", packageName)
	if err != nil {
		return nil, err
	}

	regexString := fmt.Sprintf("'%s' \\((.*)\\) was installed", packageName)
	versionCaptureRegex, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}

	capturedVersion, err := scoopPackageManager.shellCommandService.RunShellCommand(scoopPackageManager.Name(), true,
		versionCaptureRegex, "install", packageName)
	if err != nil || capturedVersion == "" {
		if err == nil {
			err = fmt.Errorf("error installing package")
		}
		return nil, err
	}

	return NewVersion(capturedVersion), nil
}

// UpdatePackage updates the package of the given name. If a version is given, that specific version of the package is
// installed. Otherwise, the latest version is installed.
//
// It returns the version of the package that was installed.
func (scoopPackageManager *ScoopPackageManager) UpdatePackage(packageName string, version *Version) (*Version, error) {
	_, err := fmt.Fprintf(scoopPackageManager.outputWriter, "Updating package \"%s\"...\n", packageName)
	if err != nil {
		return nil, err
	}

	regexString := fmt.Sprintf("'%s' \\((.*)\\) was installed", packageName)
	versionCaptureRegex, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}

	capturedVersion, err := scoopPackageManager.shellCommandService.RunShellCommand(scoopPackageManager.Name(), true,
		versionCaptureRegex, "update", packageName)
	if err != nil || capturedVersion == "" {
		if err == nil {
			err = fmt.Errorf("error updating package")
		}
		return nil, err
	}

	return NewVersion(capturedVersion), nil
}

// UninstallPackage uninstalls the package of the given name.
func (scoopPackageManager *ScoopPackageManager) UninstallPackage(packageName string) error {
	_, err := fmt.Fprintf(scoopPackageManager.outputWriter, "Uninstalling package \"%s\"...\n", packageName)
	if err != nil {
		return err
	}

	// Scoop doesn't return non-zero exit codes, so we have to check the output to see if the operation was successful.
	regexString := fmt.Sprintf("('%s' was uninstalled)", packageName)
	successRegex, err := regexp.Compile(regexString)
	if err != nil {
		return err
	}

	capturedSuccess, err := scoopPackageManager.shellCommandService.RunShellCommand(scoopPackageManager.Name(), true,
		successRegex, "uninstall", packageName)
	if err != nil || capturedSuccess == "" {
		if err == nil {
			err = fmt.Errorf("error uninstalling package")
		}
		return err
	}

	return nil
}
