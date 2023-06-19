package packagemanagers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPackagemanagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Packagemanagers Suite")
}
