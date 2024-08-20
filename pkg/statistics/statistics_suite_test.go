package statistics_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStatistics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Statistics Suite")
}
