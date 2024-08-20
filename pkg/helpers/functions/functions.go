package functions

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/helpers/functions")

// PickRandomDate picks a random date.
func PickRandomDate(ctx context.Context) string {
	_, span := tracer.Start(ctx, "PickRandomDate")
	defer span.End()

	year := rand.Intn(2100-1900) + 1900
	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1
	return fmt.Sprintf("%0.4d-%0.2d-%0.2d", year, month, day)
}

// PickRandomStringFromSlice picks a random element from the given slice.
func PickRandomStringFromSlice(ctx context.Context, slice *[]string) string {
	_, span := tracer.Start(ctx, "PickRandomStringFromSlice")
	defer span.End()

	if len(*slice) == 0 {
		return ""
	}
	randIndex := rand.Intn(len(*slice))
	return (*slice)[randIndex]
}

// PickRandomSliceFromSlice picks a random slice from the given slice.
func PickRandomSliceFromSlice(ctx context.Context, slice *[][]string) []string {
	_, span := tracer.Start(ctx, "PickRandomSliceFromSlice")
	defer span.End()

	if len(*slice) == 0 {
		return []string{}
	}
	randIndex := rand.Intn(len(*slice))
	return (*slice)[randIndex]
}

// PickRandomYear picks a random year.
func PickRandomYear(ctx context.Context) string {
	_, span := tracer.Start(ctx, "PickRandomYear")
	defer span.End()

	year, _, _ := time.Now().Date()
	return fmt.Sprintf("%d", rand.Intn(year-1900)+1900)
}

// RandomBase64String returns a random base64 string.
func RandomBase64String(ctx context.Context) string {
	_, span := tracer.Start(ctx, "RandomBase64String")
	defer span.End()

	length := rand.Intn(500-100) + 100
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return base64.StdEncoding.EncodeToString(b)
}

// RecalculateProbabilityWithUncertainity recalculates the probability with the given uncertainty.
func RecalculateProbabilityWithUncertainity(ctx context.Context, baseProbability float64, uncertainty float64, definePrefix int) float64 {
	_, span := tracer.Start(ctx, "RecalculateProbabilityWithUncertainity")
	defer span.End()

	var prefix int
	// if definePrefix is 0, then the probability will be calculated with the baseProbability
	// this is mostly for testing & coverage purposes, in production definePrefix should be set to 0
	// When you set the definePrefix, even numbers will decrease the probability and odd numbers will increase the probability
	if definePrefix == 0 {
		prefix = rand.Intn(100)
	} else {
		prefix = definePrefix
	}
	if prefix%2 == 0 {
		return baseProbability - rand.Float64()*uncertainty
	}
	return baseProbability + rand.Float64()*uncertainty
}

// SleepWithContext sleeps for the given duration or until the context is done.
func SleepWithContext(ctx context.Context, duration time.Duration) {
	// No need to trace this function as it has a fixed runtime.

	t := time.NewTimer(duration)
	select {
	case <-ctx.Done():
		t.Stop()
		fmt.Println("Context time interrupted")
	case <-t.C:
	}
}
