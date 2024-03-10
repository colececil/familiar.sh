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
	outputWriter                 io.Writer
	operatingSystemService       *system.OperatingSystemService
	createShellCommandFunc       system.CreateShellCommandFunc
	createShellCommandRunnerFunc system.CreateShellCommandRunnerFunc
}

// NewScoopPackageManager returns a new instance of ScoopPackageManager.
func NewScoopPackageManager(outputWriter io.Writer, operatingSystemService *system.OperatingSystemService,
	createShellCommandFunc system.CreateShellCommandFunc,
	createShellCommandRunnerFunc system.CreateShellCommandRunnerFunc) *ScoopPackageManager {

	return &ScoopPackageManager{
		operatingSystemService:       operatingSystemService,
		outputWriter:                 outputWriter,
		createShellCommandFunc:       createShellCommandFunc,
		createShellCommandRunnerFunc: createShellCommandRunnerFunc,
	}
}

// Name returns the name of the package manager.
func (s *ScoopPackageManager) Name() string {
	return "scoop"
}

// Order implements PackageManager.Order by returning 1.
func (s *ScoopPackageManager) Order() int {
	return 1
}

// IsSupported returns whether the package manager is supported on the current machine.
func (s *ScoopPackageManager) IsSupported() bool {
	return s.operatingSystemService.IsWindows()
}

// IsInstalled returns true if the package manager is installed.
func (s *ScoopPackageManager) IsInstalled() (bool, error) {
	_, err := fmt.Fprintf(s.outputWriter, "Checking if package manager \"%s\" is installed...\n", s.Name())
	if err != nil {
		return false, err
	}

	shellCommandRunner := s.createShellCommandRunnerFunc(s.createShellCommandFunc, nil, s.Name(), "--version")
	_, err = shellCommandRunner.Run(nil)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// Install installs the package manager.
func (s *ScoopPackageManager) Install() error {
	_, err := fmt.Fprintf(s.outputWriter, "Installing package manager \"%s\"...\n",
		s.Name())
	if err != nil {
		return err
	}

	shellCommandRunner := s.createShellCommandRunnerFunc(s.createShellCommandFunc, s.outputWriter, "powershell",
		"irm get.scoop.sh | iex")
	_, err = shellCommandRunner.Run(nil)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the package manager.
func (s *ScoopPackageManager) Update() error {
	_, err := fmt.Fprintf(s.outputWriter, "Updating package manager \"%s\"...\n", s.Name())
	if err != nil {
		return err
	}

	// Scoop doesn't return non-zero exit codes, so we have to check the output to see if the operation was successful.
	regexString := "(Scoop was updated)"
	successRegex, err := regexp.Compile(regexString)
	if err != nil {
		return err
	}

	shellCommandRunner := s.createShellCommandRunnerFunc(s.createShellCommandFunc, s.outputWriter, s.Name(), "update")
	capturedSuccess, err := shellCommandRunner.Run(successRegex)
	if err != nil || capturedSuccess == "" {
		if err == nil {
			err = fmt.Errorf("error updating package manager \"%s\"", s.Name())
		}
		return err
	}

	return nil
}

// Uninstall uninstalls the package manager.
func (s *ScoopPackageManager) Uninstall() error {
	_, err := fmt.Fprintf(s.outputWriter, "Uninstalling package manager \"%s\"...\n", s.Name())
	if err != nil {
		return err
	}

	// Scoop doesn't return non-zero exit codes, so we have to check the output to see if the operation was successful.
	regexString := fmt.Sprintf("('%s' was uninstalled)", s.Name())
	successRegex, err := regexp.Compile(regexString)
	if err != nil {
		return err
	}

	shellCommandRunner := s.createShellCommandRunnerFunc(s.createShellCommandFunc, s.outputWriter, s.Name(),
		"uninstall", s.Name())
	capturedSuccess, err := shellCommandRunner.Run(successRegex)
	if err != nil || capturedSuccess == "" {
		if err == nil {
			err = fmt.Errorf("error uninstalling package manager \"%s\"", s.Name())
		}
		return err
	}

	return nil
}

// InstalledPackages returns a slice containing information about all packages that are installed.
func (s *ScoopPackageManager) InstalledPackages() ([]*Package, error) {
	_, err := fmt.Fprintf(s.outputWriter, "Getting installed package information from package manager \"%s\"...\n",
		s.Name())
	if err != nil {
		return nil, err
	}

	jsonCaptureRegex, err := regexp.Compile("(?s)(.*)")
	if err != nil {
		return nil, err
	}

	shellCommandRunner := s.createShellCommandRunnerFunc(s.createShellCommandFunc, nil, s.Name(), "export")
	capturedJson, err := shellCommandRunner.Run(jsonCaptureRegex)
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

	shellCommandRunner = s.createShellCommandRunnerFunc(s.createShellCommandFunc, nil, s.Name(), "status")
	capturedPackages, err := shellCommandRunner.Run(packagesCaptureRegex)
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
func (s *ScoopPackageManager) InstallPackage(packageName string, version *Version) (*Version, error) {
	_, err := fmt.Fprintf(s.outputWriter, "Installing package \"%s\"...\n", packageName)
	if err != nil {
		return nil, err
	}

	regexString := fmt.Sprintf("'%s' \\((.*)\\) was installed", packageName)
	versionCaptureRegex, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}

	shellCommandRunner := s.createShellCommandRunnerFunc(s.createShellCommandFunc, s.outputWriter, s.Name(), "install",
		packageName)
	capturedVersion, err := shellCommandRunner.Run(versionCaptureRegex)
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
func (s *ScoopPackageManager) UpdatePackage(packageName string, version *Version) (*Version, error) {
	_, err := fmt.Fprintf(s.outputWriter, "Updating package \"%s\"...\n", packageName)
	if err != nil {
		return nil, err
	}

	regexString := fmt.Sprintf("'%s' \\((.*)\\) was installed", packageName)
	versionCaptureRegex, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}

	shellCommandRunner := s.createShellCommandRunnerFunc(s.createShellCommandFunc, s.outputWriter, s.Name(), "update",
		packageName)
	capturedVersion, err := shellCommandRunner.Run(versionCaptureRegex)
	if err != nil || capturedVersion == "" {
		if err == nil {
			err = fmt.Errorf("error updating package")
		}
		return nil, err
	}

	return NewVersion(capturedVersion), nil
}

// UninstallPackage uninstalls the package of the given name.
func (s *ScoopPackageManager) UninstallPackage(packageName string) error {
	_, err := fmt.Fprintf(s.outputWriter, "Uninstalling package \"%s\"...\n", packageName)
	if err != nil {
		return err
	}

	// Scoop doesn't return non-zero exit codes, so we have to check the output to see if the operation was successful.
	regexString := fmt.Sprintf("('%s' was uninstalled)", packageName)
	successRegex, err := regexp.Compile(regexString)
	if err != nil {
		return err
	}

	shellCommandRunner := s.createShellCommandRunnerFunc(s.createShellCommandFunc, s.outputWriter, s.Name(),
		"uninstall", packageName)
	capturedSuccess, err := shellCommandRunner.Run(successRegex)
	if err != nil || capturedSuccess == "" {
		if err == nil {
			err = fmt.Errorf("error uninstalling package")
		}
		return err
	}

	return nil
}
