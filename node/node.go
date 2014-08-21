package node

import "fmt"

type Node struct {
	Name string // Unique identifier and node fqdn
	Role string // Role this node is assigned to
}

// NewNode creates a new node given a name and role, neither can be empty.
func NewNode(name string, role string) (*Node, error) {
	if name == "" {
		return nil, fmt.Errorf("empty name")
	}
	if role == "" {
		return nil, fmt.Errorf("empty role")
	}
	return &Node{name, role}, nil
}

// NodeManager manages a list of nodes in memory.
type NodeManager struct {
	nodes []*Node
}

// NewNodeManager returns an empty NodeManager.
func NewNodeManager() *NodeManager {
	return &NodeManager{}
}

// Save saves the given Node in the NodeManager.
func (m *NodeManager) Save(node *Node) error {
	for i, n := range m.nodes {
		if n.Name == node.Name {
			m.nodes[i] = cloneNode(node)
			return nil
		}
	}

	// If node.Name is not found append
	m.nodes = append(m.nodes, cloneNode(node))
	return nil
}

// cloneNode creates and returns a deep copy of the given Node.
func cloneNode(n *Node) *Node {
	c := *n
	return &c
}

// All returns the list of all the Nodes in the NodeManager.
func (m *NodeManager) All() []*Node {
	return m.nodes
}

// Find returns the Node with the given name in the NodeManager and a boolean
// indicating if the id was found.
func (m *NodeManager) Find(Name string) (*Node, bool) {
	for _, n := range m.nodes {
		if n.Name == Name {
			return n, true
		}
	}
	return nil, false
}
