package functions_test

import (
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFunctions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Functions Suite")
}

var _ = Describe("Functions", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("PickRandomDate", func() {
		It("should return a random date", func() {
			randomDate := functions.PickRandomDate(ctx)
			Expect(randomDate).ToNot(BeEmpty())
			Expect(randomDate).To(MatchRegexp(`\d{4}-\d{2}-\d{2}`))
		})
	})

	Context("PickRandomStringFromSlice", func() {
		It("should return a random string from the slice", func() {
			slice := []string{"a", "b", "c"}
			randomString := functions.PickRandomStringFromSlice(ctx, &slice)
			Expect(randomString).ToNot(BeEmpty())
			Expect(slice).To(ContainElement(randomString))
		})

		It("should return an empty string if the slice is empty", func() {
			slice := []string{}
			randomString := functions.PickRandomStringFromSlice(ctx, &slice)
			Expect(randomString).To(BeEmpty())
		})

		It("should return an empty string if the slice is nil", func() {
			var slice []string
			randomString := functions.PickRandomStringFromSlice(ctx, &slice)
			Expect(randomString).To(BeEmpty())
		})

		It("should return the first item if the slice has only one item", func() {
			slice := []string{"a"}
			randomString := functions.PickRandomStringFromSlice(ctx, &slice)
			Expect(randomString).To(Equal("a"))
		})
	})

	Context("PickRandomSliceFromSlice", func() {
		It("should return a random slice from the slice", func() {
			slice := [][]string{{"a", "b"}, {"c", "d"}}
			randomSlice := functions.PickRandomSliceFromSlice(ctx, &slice)
			Expect(randomSlice).ToNot(BeEmpty())
			Expect(slice).To(ContainElement(randomSlice))
		})

		It("should return an empty slice if the slice is empty", func() {
			slice := [][]string{}
			randomSlice := functions.PickRandomSliceFromSlice(ctx, &slice)
			Expect(randomSlice).To(BeEmpty())
		})

		It("should return an empty slice if the slice is nil", func() {
			var slice [][]string
			randomSlice := functions.PickRandomSliceFromSlice(ctx, &slice)
			Expect(randomSlice).To(BeEmpty())
		})

		It("should return the first item if the slice has only one item", func() {
			slice := [][]string{{"a"}}
			randomSlice := functions.PickRandomSliceFromSlice(ctx, &slice)
			Expect(randomSlice).To(Equal([]string{"a"}))
		})
	})

	Context("PickRandomYear", func() {
		It("should return a random year", func() {
			randomYear := functions.PickRandomYear(ctx)
			Expect(randomYear).ToNot(BeEmpty())
			Expect(randomYear).To(MatchRegexp(`\d{4}`))
		})
	})

	Context("RandomBase64String", func() {
		It("should return a random base64 string", func() {
			randomBase64String := functions.RandomBase64String(ctx)
			Expect(randomBase64String).ToNot(BeEmpty())
		})

		It("should return a different string on each call", func() {
			randomBase64String1 := functions.RandomBase64String(ctx)
			randomBase64String2 := functions.RandomBase64String(ctx)
			Expect(randomBase64String1).ToNot(Equal(randomBase64String2))
		})
	})

	Context("RecalculateProbabilityWithUncertainity", func() {
		It("should recalculate the probability with the given uncertainty", func() {
			baseProbability := 0.5
			uncertainty := 0.1
			recalculatedProbability := functions.RecalculateProbabilityWithUncertainity(ctx, baseProbability, uncertainty, 0)
			Expect(recalculatedProbability).To(BeNumerically("~", baseProbability, uncertainty))
		})

		It("should return the base probability if the uncertainty is 0", func() {
			baseProbability := 0.5
			recalculatedProbability := functions.RecalculateProbabilityWithUncertainity(ctx, baseProbability, 0, 0)
			Expect(recalculatedProbability).To(Equal(baseProbability))
		})

		It("should return a value lower than the base probability if the definePrefix is negative", func() {
			baseProbability := 0.5
			recalculatedProbability := functions.RecalculateProbabilityWithUncertainity(ctx, baseProbability, 0.1, 10)
			Expect(recalculatedProbability).To(BeNumerically("<", baseProbability))
		})

		It("should return a value higher than the base probability if the definePrefix is positive", func() {
			baseProbability := 0.5
			recalculatedProbability := functions.RecalculateProbabilityWithUncertainity(ctx, baseProbability, 0.1, 5)
			Expect(recalculatedProbability).To(BeNumerically(">", baseProbability))
		})
	})

	Context("SleepWithContext", func() {
		It("should sleep for the given duration", func() {
			start := time.Now()
			duration := 100 * time.Millisecond
			ctx, cancel := context.WithTimeout(ctx, duration)
			defer cancel()
			functions.SleepWithContext(ctx, duration)
			Expect(time.Since(start)).To(BeNumerically("~", duration, 10*time.Millisecond))
		})

		It("should return immediately if the context is done", func() {
			start := time.Now()
			duration := 100 * time.Millisecond
			ctx, cancel := context.WithTimeout(ctx, 0)
			defer cancel()
			functions.SleepWithContext(ctx, duration)
			Expect(time.Since(start)).To(BeNumerically("<", duration))
		})
	})
})
