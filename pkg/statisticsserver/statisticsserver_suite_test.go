package statisticsserver_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStatisticsserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Statisticsserver Suite")
}
