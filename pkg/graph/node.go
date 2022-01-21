package graph

import (
	"context"
	"encoding/json"
	"math/rand"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	Exit = "exit"
)

var (
	Directions = []string{
		"left",
		"right",
		"forward",
		"upstairs",
		"downstairs",
	}

	Rooms = []string{
		"dragon",
		"troll",
		"dark knight",
		"ninja",
		"zombie",
	}
)

type Node map[string]json.RawMessage

// FindExit recursively transforms JSON message into Path to exit.
// depth argument is used to allocate final slice.
func (n Node) FindExit(ctx context.Context, depth uint64) (*Path, error) {
	for key, value := range n {
		// timeout management
		select {
		case _ = <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var child Node

		if err := json.Unmarshal(value, &child); err != nil {
			// ignore error, fallback on unmarshal string
		} else {
			g, err := child.FindExit(ctx, depth+1)
			if err != nil {
				return nil, err
			}

			// exit path found, return
			if g.Directions != nil {
				g.Directions[depth] = key

				return g, nil
			}

			// child node detected, don't try unmarshal string
			continue
		}

		var s string

		if err := json.Unmarshal(value, &s); err != nil {
			return nil, err
		}

		// exit found, return
		if s == Exit {
			d := make([]string, depth+1)
			d[depth] = key

			return &Path{
				Directions: d,
			}, nil
		}
	}

	return &Path{
		Directions: nil,
	}, nil
}

type ConfigGeneration struct {
	Width  int
	Height int

	// Function to create random room names
	RoomFunc func() string

	// Unique exit management
	ExitFunc    func() bool // ExitFunc returns if a node is an exit or not and ensures only one exit is created
	ExitCreated *struct{}   // ExitCreated is nil if no exit is set yet
}

func (n *Node) Generate(ctx context.Context, cfg *ConfigGeneration, eg *errgroup.Group, depth int) error {
	*n = make(Node)

	// Lock mutex for conccurent map write
	// NOT IDEAL SOLUTION, sync.Map is designed for this but it would slow down in many cases
	// and enforces two different Node models.
	var wm sync.Mutex

	// usage of weak random generator is ok here
	nchildren := rand.Intn(cfg.Width + 1) // nolint: gosec

	// shuffle directions to have a bit of everything at each new node
	directions := make([]string, len(Directions))
	copy(directions, Directions)
	rand.Shuffle(len(directions), func(i, j int) { directions[i], directions[j] = directions[j], directions[i] })

	// children down value depth
	depth--

	for i := 0; i < nchildren; i++ {
		// timeout management
		select {
		case _ = <-ctx.Done():
			return ctx.Err()
		default:
		}

		// if max depth is reached, create a leaf with exit roll
		if depth == 0 {
			// roll if no-go room is exit or not
			if cfg.ExitCreated == nil && cfg.ExitFunc() {
				cfg.ExitCreated = &struct{}{}
				(*n)[directions[i]] = json.RawMessage(`"` + Exit + `"`)
			} else {
				(*n)[directions[i]] = json.RawMessage(`"` + cfg.RoomFunc() + `"`)
			}
		} else {
			i := i

			eg.Go(func() error {
				var child Node

				var branchEG errgroup.Group

				// create a new sub graph
				if err := child.Generate(ctx, cfg, &branchEG, depth); err != nil {
					return err
				}

				if err := branchEG.Wait(); err != nil {
					return err
				}

				// sub graph is actually a leaf
				if len(child) == 0 {
					wm.Lock()
					(*n)[directions[i]] = json.RawMessage(`"` + cfg.RoomFunc() + `"`)
					wm.Unlock()
				} else {
					// set subgraph as json raw message
					raw, err := json.Marshal(child)
					if err != nil {
						return err
					}

					wm.Lock()
					(*n)[directions[i]] = json.RawMessage(raw)
					wm.Unlock()
				}

				return nil
			})
		}
	}

	return nil
}
