package graph

import (
	"context"
	"encoding/json"
)

const (
	Exit = "exit"
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
