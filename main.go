// Completelty written by Claude Opus 4.8 Just for testing purposes
package main

import (
	"fmt"
	"time"

	"github.com/notblankz/un-ramen/dag"
)

// task returns a stub run func: announce, simulate work, announce done.
func task(name string, work time.Duration) func() error {
	return func() error {
		fmt.Printf("  → starting  %s\n", name)
		time.Sleep(work)
		fmt.Printf("  ✓ finished  %s\n", name)
		return nil
	}
}

func main() {
	g := dag.NewGraph()

	// the research pipeline:
	//   fetch ──> summarize ──┐
	//         └─> extract  ───┴──> synthesize ──> critique
	g.AddNode("fetch", task("fetch", 400*time.Millisecond))
	g.AddNode("summarize", task("summarize", 1*time.Second))
	g.AddNode("extract", task("extract", 1*time.Second))
	g.AddNode("synthesize", task("synthesize", 500*time.Millisecond))
	g.AddNode("critique", task("critique", 300*time.Millisecond))

	edges := [][2]string{
		{"fetch", "summarize"},
		{"fetch", "extract"},
		{"summarize", "synthesize"},
		{"extract", "synthesize"},
		{"synthesize", "critique"},
	}
	for _, e := range edges {
		if err := g.AddEdge(e[0], e[1]); err != nil {
			fmt.Println("build error:", err)
			return
		}
	}

	g.Print()
	fmt.Println()

	start := time.Now()
	if err := g.RunParallel(); err != nil {
		fmt.Println("run error:", err)
		return
	}
	fmt.Printf("\ntotal: %s\n", time.Since(start))
}
