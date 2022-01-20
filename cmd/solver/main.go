package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/elojah/pxlmtc-test/pkg/graph"
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

	// Read JSON input
	var n graph.Node

	if err := json.NewDecoder(os.Stdin).Decode(&n); err != nil {
		logger.Error().Err(err).Msg("failed to JSON decode input")

		return
	}

	// Find and build exit path
	g, err := n.FindExit(ctx, 0)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse input")

		return
	}

	// No exit found
	if g.Directions == nil {
		fmt.Println("Sorry")

		return
	}

	// Display result
	raw, err := json.MarshalIndent(g.Directions, "", " ")
	if err != nil {
		logger.Error().Err(err).Msg("failed to display path")

		return
	}

	fmt.Println(string(raw))
}
