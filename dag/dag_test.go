package dag

import (
	"reflect"
	"testing"
)

// diamond builds A→B, A→C, B→D, C→D
func diamond(t *testing.T) *Graph {
	t.Helper()
	g := NewGraph()
	for _, id := range []string{"A", "B", "C", "D"} {
		g.AddNode(id, nil)
	}
	for _, e := range [][2]string{{"A", "B"}, {"A", "C"}, {"B", "D"}, {"C", "D"}} {
		if err := g.AddEdge(e[0], e[1]); err != nil {
			t.Fatalf("AddEdge(%q, %q): %v", e[0], e[1], err)
		}
	}
	return g
}

// assertValidTopoOrder checks the ordering *property* rather than one exact
// permutation: every node appears exactly once, and no node precedes something
// it depends on.
func assertValidTopoOrder(t *testing.T, g *Graph, order []string) {
	t.Helper()

	if len(order) != len(g.nodes) {
		t.Fatalf("order has %d nodes, graph has %d", len(order), len(g.nodes))
	}

	pos := make(map[string]int, len(order))
	for i, id := range order {
		if _, dup := pos[id]; dup {
			t.Fatalf("node %q appears more than once", id)
		}
		if _, ok := g.nodes[id]; !ok {
			t.Fatalf("order contains unknown node %q", id)
		}
		pos[id] = i
	}

	for fromID, children := range g.edges {
		for _, toID := range children {
			if pos[fromID] > pos[toID] {
				t.Errorf("edge %s→%s violated: %s at index %d, %s at index %d",
					fromID, toID, fromID, pos[fromID], toID, pos[toID])
			}
		}
	}
}

func TestTopoSortDiamond(t *testing.T) {
	g := diamond(t)

	order, err := g.TopoSort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertValidTopoOrder(t, g, order)
}

func TestTopoSortDetectsCycle(t *testing.T) {
	g := NewGraph()
	for _, id := range []string{"A", "B", "C"} {
		g.AddNode(id, nil)
	}
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "A") // closes the loop

	if _, err := g.TopoSort(); err == nil {
		t.Fatal("expected a cycle error, got nil")
	}
}

// multiRoot has TWO zero-in-degree nodes, so the seeding loop's map range
// can produce different orders — this is what actually exercises the sort.
func multiRoot(t *testing.T) *Graph {
	t.Helper()
	g := NewGraph()
	for _, id := range []string{"A", "B", "C", "D"} {
		g.AddNode(id, nil)
	}
	for _, e := range [][2]string{{"A", "C"}, {"B", "C"}, {"C", "D"}} {
		if err := g.AddEdge(e[0], e[1]); err != nil {
			t.Fatalf("AddEdge(%q, %q): %v", e[0], e[1], err)
		}
	}
	return g
}

func TestTopoSortIsDeterministic(t *testing.T) {
	first, err := multiRoot(t).TopoSort()
	if err != nil {
		t.Fatal(err)
	}

	for i := range 20 {
		got, err := multiRoot(t).TopoSort()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(first, got) {
			t.Fatalf("run %d differs:\nfirst: %v\ngot:   %v", i, first, got)
		}
	}
}

func TestAddEdgeUnknownNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("A", nil)

	if err := g.AddEdge("A", "ghost"); err == nil {
		t.Error("expected error for unknown target, got nil")
	}
	if err := g.AddEdge("ghost", "A"); err == nil {
		t.Error("expected error for unknown source, got nil")
	}
}

func TestTopoSortEmptyGraph(t *testing.T) {
	order, err := NewGraph().TopoSort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 0 {
		t.Fatalf("expected empty order, got %v", order)
	}
}
