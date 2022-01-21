package graph_test

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"testing"

	"github.com/elojah/pxlmtc-test/pkg/graph"
	"golang.org/x/sync/errgroup"
)

func TestFindExit(t *testing.T) {
	data := []struct {
		n          graph.Node
		inputCtx   context.Context
		inputDepth uint64
		outputPath graph.Path
		outputErr  error
	}{
		{
			n: graph.Node{
				"forward": json.RawMessage(`"exit"`),
			},
			inputCtx:   context.Background(),
			inputDepth: 0,
			outputPath: graph.Path{Directions: []string{"forward"}},
			outputErr:  nil,
		},
		{
			n: graph.Node{
				"forward": json.RawMessage(`"tiger"`),
				"left":    json.RawMessage(`{"forward": {"upstairs": "exit"}, "left": "dragon"}`),
				"right":   json.RawMessage(`{"forward":"dead end"}}`),
			},
			inputCtx:   context.Background(),
			inputDepth: 0,
			outputPath: graph.Path{Directions: []string{"left", "forward", "upstairs"}},
			outputErr:  nil,
		},
		{
			n: graph.Node{
				"forward": json.RawMessage(`"tiger"`),
			},
			inputCtx:   context.Background(),
			inputDepth: 0,
			outputPath: graph.Path{Directions: []string{}},
			outputErr:  nil,
		},
	}

	for _, d := range data {
		p, err := d.n.FindExit(d.inputCtx, d.inputDepth)

		if !errors.Is(err, d.outputErr) {
			t.Error("wrong error")

			continue
		}

		if len(p.Directions) != len(d.outputPath.Directions) {
			t.Error("path has bad length")

			continue
		}

		for i, dir := range p.Directions {
			if dir != d.outputPath.Directions[i] {
				t.Error("path has bad value")

				break
			}
		}
	}
}

func BenchmarkFindExit(b *testing.B) {
	n := graph.Node{
		"forward": json.RawMessage(`"tiger"`),
		"left":    json.RawMessage(`{"forward": {"upstairs": "exit"}, "left": "dragon"}`),
		"right":   json.RawMessage(`{"forward":"dead end"}}`),
	}

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n.FindExit(ctx, 0)
	}
}

func TestGenerate(t *testing.T) {
	var once sync.Once

	data := []struct {
		n          graph.Node
		inputCtx   context.Context
		inputCfg   graph.ConfigGeneration
		inputDepth int
		outputErr  error
	}{
		{
			n:          graph.Node{},
			inputCtx:   context.Background(),
			inputCfg:   graph.ConfigGeneration{},
			inputDepth: 0,
			outputErr:  nil,
		},
		// Classic test with same exit func than main
		{
			n:        graph.Node{},
			inputCtx: context.Background(),
			inputCfg: graph.ConfigGeneration{
				Width:    3,
				Height:   5,
				RoomFunc: func() string { return "test_room" },
				ExitFunc: func() bool {
					var result bool

					// usage of weak random generator is ok here
					r := rand.Intn(100) // nolint: gosec
					if r <= 50 {
						once.Do(func() {
							result = true
						})
					}

					return result
				},
			},
			inputDepth: 5,
			outputErr:  nil,
		},
		// Test possible multiple exit
		{
			n:        graph.Node{},
			inputCtx: context.Background(),
			inputCfg: graph.ConfigGeneration{
				Width:    3,
				Height:   5,
				RoomFunc: func() string { return "test_room" },
				ExitFunc: func() bool {
					return rand.Intn(2) == 2 // nolint: gosec
				},
			},
			inputDepth: 5,
			outputErr:  nil,
		},
	}

	for _, d := range data {
		eg := errgroup.Group{}
		err := d.n.Generate(d.inputCtx, &d.inputCfg, &eg, d.inputDepth)

		if !errors.Is(err, d.outputErr) {
			t.Error("wrong error")

			continue
		}

		err = eg.Wait()

		if !errors.Is(err, d.outputErr) {
			t.Error("wrong error")

			continue
		}
	}
}

// Benchmark on generate is hard to evaluate due to random nature of function.
// func BenchmarkGenerate(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 	}
// }
