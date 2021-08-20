package capacity_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCapacity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Capacity Suite")
}
