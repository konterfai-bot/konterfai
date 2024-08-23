package hallucinator_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHallucinator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hallucinator Suite")
}
