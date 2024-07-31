package functions

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

// PickRandomDate picks a random date.
func PickRandomDate() string {
	year := rand.Intn(2100-1900) + 1900
	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1
	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

// PickRandomFromSlice picks a random element from the given slice.
func PickRandomFromSlice(slice *[]string) string {
	if len(*slice) == 0 {
		return ""
	}
	randIndex := rand.Intn(len(*slice))
	return (*slice)[randIndex]
}

// PickRandomYear picks a random year.
func PickRandomYear() string {
	year, _, _ := time.Now().Date()
	return fmt.Sprintf("%d", rand.Intn(year-1900)+1900)
}

// RandomBase64String returns a random base64 string.
func RandomBase64String() string {
	length := rand.Intn(500-100) + 100
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return base64.StdEncoding.EncodeToString(b)
}

// RecalculateProbabilityWithUncertainity recalculates the probability with the given uncertainty.
func RecalculateProbabilityWithUncertainity(baseProbability float64, uncertainty float64) float64 {
	prefix := rand.Intn(100)
	if prefix%2 == 0 {
		return baseProbability - rand.Float64()*uncertainty
	}
	return baseProbability + rand.Float64()*uncertainty
}

// SleepWithContext sleeps for the given duration or until the context is done.
func SleepWithContext(ctx context.Context, duration time.Duration) {
	t := time.NewTimer(duration)
	select {
	case <-ctx.Done():
		t.Stop()
		fmt.Println("Context time interrupted")
	case <-t.C:
	}
}
