package main

import (
	"context"
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	version string
	prog    string
)

func main() {
	// Initialize context for timeout
	ctx := context.Background()

	// Initialize logger
	zerolog.TimeFieldFormat = ""

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger = logger.With().Str("version", version).Str("exe", prog).Logger()

	// Read timeout flag
	to := flag.String("timeout", "0", "timeout computation value (default: 0)")
	width := flag.Uint64("width", 3, "maximum number of branches out of room")
	height := flag.Uint64("width", 10, "maximum depth to exit")

	flag.Parse()

	if to != nil {
		d, err := time.ParseDuration(*to)
		if err != nil {
			logger.Error().Err(err).Msg("failed to parse timeout duration")

			return
		}

		if d > 0 {
			var cancel context.CancelFunc

			ctx, cancel = context.WithTimeout(ctx, d)
			defer func() {
				cancel()
			}()
		}
	}

	_ = ctx

	rand.Seed(time.Now().UnixNano())
}
