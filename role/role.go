package role

import "fmt"

type Role struct {
	Name string // Unique identifier and puppet class name
}

// NewRole creates a new role given a name, that can't be empty.
func NewRole(name string) (*Role, error) {
	if name == "" {
		return nil, fmt.Errorf("empty name")
	}
	return &Role{name}, nil
}

// RoleManager manages a list of roles in memory.
type RoleManager struct {
	roles []*Role
}

// NewRoleManager returns an empty RoleManager.
func NewRoleManager() *RoleManager {
	return &RoleManager{}
}

// Save saves the given Role in the RoleManager.
func (m *RoleManager) Save(role *Role) error {
	for i, r := range m.roles {
		if r.Name == role.Name {
			m.roles[i] = cloneRole(role)
			return nil
		}
	}

	// If role.Name is not found append
	m.roles = append(m.roles, cloneRole(role))
	return nil
}

// cloneRole creates and returns a deep copy of the given Role.
func cloneRole(r *Role) *Role {
	c := *r
	return &c
}

// All returns the list of all the Roles in the RoleManager.
func (m *RoleManager) All() []*Role {
	return m.roles
}

// Find returns the Role with the given name in the RoleManager and a boolean
// indicating if the id was found.
func (m *RoleManager) Find(Name string) (*Role, bool) {
	for _, r := range m.roles {
		if r.Name == Name {
			return r, true
		}
	}
	return nil, false
}
