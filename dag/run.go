package dag

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

func (g *Graph) Run() error {
	// This function calls the TopoSort and then executes each node's run function
	sorted, err := g.TopoSort()
	if err != nil {
		return err
	}

	for _, nodeID := range sorted {
		node := g.nodes[nodeID]
		if node.run == nil {
			continue
		}
		err = node.run()
		if err != nil {
			return fmt.Errorf("dag: task %q failed: %w", nodeID, err)
		}
	}

	return err
}

func (g *Graph) RunParallel() error {

	levels, err := g.Levels()
	if err != nil {
		return err
	}

	for _, level := range levels {
		var grp errgroup.Group
		for _, nodeID := range level {
			grp.Go(func() error {
				node := g.nodes[nodeID]
				if node.run == nil {
					return nil
				}

				return node.run()
			})
		}
		if err := grp.Wait(); err != nil {
			return err
		}
	}

	return nil
}
