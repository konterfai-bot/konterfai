package statisticsserver_test

import (
	"sort"

	"codeberg.org/konterfai/konterfai/pkg/statisticsserver"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestData", func() {
	var rds statisticsserver.RequestDataSlice
	BeforeEach(func() {
		rds = statisticsserver.RequestDataSlice{
			{
				Identifier:          "test3",
				Count:               3,
				Size:                "3MiB",
				IsRobotsTxtViolator: "ignored",
			},
			{
				Identifier:          "test2",
				Count:               2,
				Size:                "2MiB",
				IsRobotsTxtViolator: "yes",
			},
			{
				Identifier:          "test1",
				Count:               1,
				Size:                "1MiB",
				IsRobotsTxtViolator: "no",
			},
		}
	})

	Context("Len", func() {
		It("should return the length of the slice", func() {
			Expect(rds.Len()).To(Equal(3))
		})
	})

	Context("Swap", func() {
		It("should swap the elements at the given indices", func() {
			rds.Swap(0, 1)
			Expect(rds[0].Identifier).To(Equal("test2"))
			Expect(rds[1].Identifier).To(Equal("test3"))
			Expect(rds.Len()).To(Equal(3))
		})
	})

	Context("Less", func() {
		It("should compare the elements at the given indices", func() {
			Expect(rds.Less(0, 1)).To(BeFalse())
			Expect(rds.Less(1, 0)).To(BeTrue())
			Expect(rds.Len()).To(Equal(3))
		})
	})

	Context("Sort", func() {
		It("should sort the elements by identifier", func() {
			sort.Sort(rds)
			Expect(rds[0].Identifier).To(Equal("test1"))
			Expect(rds[1].Identifier).To(Equal("test2"))
			Expect(rds[2].Identifier).To(Equal("test3"))
			Expect(rds.Len()).To(Equal(3))
		})
	})
})
