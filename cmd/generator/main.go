package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/elojah/pxlmtc-test/pkg/graph"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
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
	width := flag.Int("width", 3, "maximum number of branches out of room")
	height := flag.Int("height", 10, "maximum depth to exit")
	exit := flag.Int("exit", 10, "roll percentage for leaf to be an exit")

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

	if *width > len(graph.Directions) {
		logger.Error().Msgf("width cannot exceed %d", len(graph.Directions))

		return
	}

	rand.Seed(time.Now().UnixNano())

	var once sync.Once

	cfg := graph.ConfigGeneration{
		Width:  *width,
		Height: *height,

		RoomFunc: func() string {
			return graph.Rooms[rand.Intn(len(graph.Rooms))] // nolint: gosec
		},

		// ExitFunc embed a sync.Once to ensure exit is generated once
		ExitFunc: func() bool {
			var result bool

			// usage of weak random generator is ok here
			r := rand.Intn(100) // nolint: gomnd, gosec
			if r <= *exit {
				once.Do(func() {
					result = true
				})
			}

			return result
		},
		ExitCreated: nil,
	}

	var n graph.Node

	var eg errgroup.Group

	if err := n.Generate(ctx, &cfg, &eg, *height); err != nil {
		logger.Error().Err(err).Msg("failed to generate graph")

		return
	}

	// Wait concurrent generation to be done
	if err := eg.Wait(); err != nil {
		logger.Error().Err(err).Msg("failed to generate graph")

		return
	}

	// Display result
	raw, err := json.MarshalIndent(n, "", " ")
	if err != nil {
		logger.Error().Err(err).Msg("failed to display path")

		return
	}

	fmt.Println(string(raw))
}
