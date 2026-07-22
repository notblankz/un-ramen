package dag

import (
	"fmt"
	"sort"
	"strings"
)

type Node struct {
	id  string
	run func() error
}

type Graph struct {
	nodes map[string]*Node
	edges map[string][]string
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*Node),
		edges: make(map[string][]string),
	}
}

// put a new node in g.nodes
func (g *Graph) AddNode(id string, run func() error) *Node {
	n := &Node{
		id:  id,
		run: run,
	}
	g.nodes[id] = n
	return n
}

// add toID into the array of g.edges[fromID]
func (g *Graph) AddEdge(fromID, toID string) error {
	// check for existence of both from and to nodes
	if _, ok := g.nodes[fromID]; !ok {
		return fmt.Errorf("dag: unknown node %q", fromID)
	}
	if _, ok := g.nodes[toID]; !ok {
		return fmt.Errorf("dag: unknown node %q", toID)
	}
	// add the edge in the
	g.edges[fromID] = append(g.edges[fromID], toID)
	return nil
}

// Written by claude for pretty printing the DAGs
func (g *Graph) Print() {
	// map iteration is randomized, so collect + sort for stable output
	ids := make([]string, 0, len(g.nodes))
	for id := range g.nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	// widest id, so all the arrows line up
	width := 0
	for _, id := range ids {
		if len(id) > width {
			width = len(id)
		}
	}

	fmt.Printf("graph · %d nodes\n\n", len(g.nodes))
	for _, id := range ids {
		children := append([]string(nil), g.edges[id]...) // copy so we don't sort the real slice
		sort.Strings(children)

		if len(children) == 0 {
			fmt.Printf("  %-*s  →  (leaf)\n", width, id)
			continue
		}
		fmt.Printf("  %-*s  →  %s\n", width, id, strings.Join(children, ", "))
	}
}

func (g *Graph) TopoSort() ([]string, error) {
	// Map to hold the indegree of all the nodes in the graph
	indeg := make(map[string]int)
	// and a queue for the sorting
	queue := make([]string, 0)
	// Final sort
	res := make([]string, 0)

	// add all nodes to indeg and set to 0
	for nodeID := range g.nodes {
		indeg[nodeID] = 0
	}

	for fromID := range g.nodes {
		// we have each node from the edges map here
		for _, toID := range g.edges[fromID] {
			// here we have each node that fromNode has an edge towards
			indeg[toID]++
		}
	}

	for nodeID, indegVal := range indeg {
		if indegVal == 0 {
			// if the indeg val of any node is 0 we add that to the queue
			queue = append(queue, nodeID)
		}
	}

	// we sort the queue to maintain reproducibility
	// between different passes of the TopoSort algo
	sort.Strings(queue)

	// We have the queue ready, now we can start the loop till queue is empty
	for len(queue) != 0 {
		// 1 Pop from queue
		// 2 Get all the nodes of the popped node and decrease their indeg by 1
		poppedNodeID := queue[0]
		queue = queue[1:]

		// decrement indeg of all the children of poppedNode
		// if after decrementing it's degree it becomes 0 we can add it to the queue
		for _, toID := range g.edges[poppedNodeID] {
			indeg[toID]--
			if indeg[toID] == 0 {
				queue = append(queue, toID)
			}
		}

		res = append(res, poppedNodeID)

	}

	// cycle check
	if len(g.nodes) != len(res) {
		return nil, fmt.Errorf("dag: cycle detected")
	}

	return res, nil
}
