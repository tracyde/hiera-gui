package node

import "testing"

func newNodeOrFatal(t *testing.T, name string, role string) *Node {
	node, err := NewNode(name, role)
	if err != nil {
		t.Fatal("new node: %v", err)
	}
	return node
}

func TestNewNode(t *testing.T) {
	name := "node.testing.com"
	role := "node::testing"
	node := newNodeOrFatal(t, name, role)
	if node.Name != name {
		t.Errorf("expected name %q, got %q", name, node.Name)
	}
	if node.Role != role {
		t.Errorf("expected role %q, got %q", role, node.Role)
	}
}

func TestNewNodeEmptyName(t *testing.T) {
	_, err := NewNode("", "node::testing")
	if err == nil {
		t.Errorf("expected 'empty name' error, got nil")
	}
}

func TestNewNodeEmptyRole(t *testing.T) {
	_, err := NewNode("node.testing.com", "")
	if err == nil {
		t.Errorf("expected 'empty role' error, got nil")
	}
}

func TestSaveNodeAndRetrieve(t *testing.T) {
	node := newNodeOrFatal(t, "node.testing.com", "node::testing")
	m := NewNodeManager()
	m.Save(node)
	all := m.All()
	if len(all) != 1 {
		t.Errorf("expected 1 node, got %v", len(all))
	}
	if *all[0] != *node {
		t.Errorf("expected %v, got %v", node, all[0])
	}
}

func TestSaveAndRetrieveTwoNodes(t *testing.T) {
	nodeTesting := newNodeOrFatal(t, "node.testing.com", "node::testing")
	nodeWWW := newNodeOrFatal(t, "www.testing.com", "node::www")
	m := NewNodeManager()
	m.Save(nodeTesting)
	m.Save(nodeWWW)
	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 nodes, got %v", len(all))
	}
	if *all[0] != *nodeTesting && *all[1] != *nodeTesting {
		t.Errorf("missing node: %v", nodeTesting)
	}
	if *all[0] != *nodeWWW && *all[1] != *nodeWWW {
		t.Errorf("missing node: %v", nodeWWW)
	}
}

func TestSaveModifyAndRetrieve(t *testing.T) {
	node := newNodeOrFatal(t, "node.testing.com", "node::testing")
	m := NewNodeManager()
	m.Save(node)
	node.Role = "node::www"
	if m.All()[0].Role == "node::www" {
		t.Errorf("saved node role was updated")
	}
}

func TestSaveTwiceAndRetrieve(t *testing.T) {
	node := newNodeOrFatal(t, "node.testing.com", "node::testing")
	m := NewNodeManager()
	m.Save(node)
	m.Save(node)
	all := m.All()
	if len(all) != 1 {
		t.Errorf("expected 1 node, got %v", len(all))
	}
	if *all[0] != *node {
		t.Errorf("expected node %v, got %v", node, all[0])
	}
}

func TestSaveAndFind(t *testing.T) {
	node := newNodeOrFatal(t, "node.testing.com", "node::testing")
	m := NewNodeManager()
	m.Save(node)
	nt, ok := m.Find(node.Name)
	if !ok {
		t.Errorf("didn't find node")
	}
	if *node != *nt {
		t.Errorf("expected %v, got %v", node, nt)
	}
}
