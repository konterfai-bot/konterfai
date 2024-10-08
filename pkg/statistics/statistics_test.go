package statistics_test

import (
	"context"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/command"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Statistics", func() {
	Context("NewStatistics", func() {
		var ctx context.Context
		var configurationInfo string
		BeforeEach(func() {
			ctx = context.Background()
			configurationInfo = "" +
				"#############################\n" +
				"# konterfAI - the anti-AI-AI #\n" +
				"#############################\n\n" +
				"Configuration:\n" +
				"\t- Address: \t\t\t\t127.0.0.1\n" +
				"\t- Port: \t\t\t\t8080\n" +
				"\t- Statistics Port: \t\t\t8081\n" +
				"\t- Generate Interval: \t\t\t5s\n" +
				"\t- Hallucination Cache Size: \t\t10\n" +
				"\t- Hallucination Prompt Word Count: \t5\n" +
				"\t- Hallucination Word Count: \t\t500\n" +
				"\t- Hallucination Request Count:  \t5\n" +
				"\t- Ollama Address: \t\t\thttp://localhost:11434\n" +
				"\t- Ollama Model: \t\t\tqwen2:0.5b\n" +
				"\t- AI Temperature: \t\t\t30\n" +
				"\t- AI Seed: \t\t\t\t0\n" +
				"\t- Hallucinator URL: \t\t\thttp://localhost:8080\n" +
				"\n"
		})

		It("should return a new statistics", func() {
			logger, _ := command.SetLogger("off", "")
			go func() {
				st := statistics.NewStatistics(ctx, logger, configurationInfo)
				Expect(st).NotTo(BeNil())
				Expect(st.ConfigurationInfo).To(Equal(configurationInfo))
				isRobotsTxt := false
				for range 10 {
					st.AppendRequest(ctx, statistics.Request{
						UserAgent:   "test",
						IPAddress:   "127.0.0.1",
						Timestamp:   time.Now(),
						IsRobotsTxt: isRobotsTxt,
						Size:        0,
					})
					isRobotsTxt = !isRobotsTxt
				}
				st.AppendRequest(ctx, statistics.Request{
					UserAgent:   "test",
					IPAddress:   "127.0.0.2",
					Timestamp:   time.Now(),
					IsRobotsTxt: false,
					Size:        0,
				})
			}()
			time.Sleep(6 * time.Second) // updates to the prometheus statistics are hardcode to 5 seconds
			ctx.Done()
		})
	})
})
