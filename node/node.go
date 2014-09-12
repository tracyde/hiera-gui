package node

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v1"
)

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
	nodeDir string
	nodes   []*Node
}

// NewNodeManager returns an empty NodeManager.
func NewNodeManager(p string) *NodeManager {
	parseDir(p)
	return &NodeManager{nodeDir: p}
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		fmt.Printf("Visited: %s\n", path)
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var n Node
		err = yaml.Unmarshal(b, &n)
		if err != nil {
			return err
		}
		fmt.Printf("node: {name: %s, role: %s}", n.Name, n.Role)
	}
	return nil
}

func parseDir(p string) error {
	return filepath.Walk(p, visit)
}

// Encode node to yaml
func (n *Node) encode(w io.Writer) (int, error) {
	y, _ := yaml.Marshal(n)
	return fmt.Fprintf(w, string(y))
}

// Write Node to disk given the node directory path
func (n *Node) write(d string) error {
	p := d + "/" + n.Name
	fmt.Println(p)
	f, err := os.Create(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		if cErr := f.Close(); cErr != nil && err == nil {
			fmt.Println(cErr)
		}
	}()
	if _, err = n.encode(f); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Save saves the given Node in the NodeManager.
func (m *NodeManager) Save(node *Node) error {
	for i, n := range m.nodes {
		if n.Name == node.Name {
			m.nodes[i] = cloneNode(node)
			return node.write(m.nodeDir)

		}
	}

	// If node.Name is not found append
	m.nodes = append(m.nodes, cloneNode(node))
	return node.write(m.nodeDir)
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
