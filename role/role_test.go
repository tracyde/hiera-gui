package role

import "testing"

func newRoleOrFatal(t *testing.T, name string) *Role {
	role, err := NewRole(name)
	if err != nil {
		t.Fatal("new role: %v", err)
	}
	return role
}

func TestNewRole(t *testing.T) {
	name := "role::testing"
	role := newRoleOrFatal(t, name)
	if role.Name != name {
		t.Errorf("expected name %q, got %q", name, role.Name)
	}
}

func TestNewRoleEmptyName(t *testing.T) {
	_, err := NewRole("")
	if err == nil {
		t.Errorf("expected 'empty name' error, got nil")
	}
}

func TestSaveRoleAndRetrieve(t *testing.T) {
	role := newRoleOrFatal(t, "role::testing")
	m := NewRoleManager()
	m.Save(role)
	all := m.All()
	if len(all) != 1 {
		t.Errorf("expected 1 role, got %v", len(all))
	}
	if *all[0] != *role {
		t.Errorf("expected %v, got %v", role, all[0])
	}
}

func TestSaveAndRetrieveTwoRoles(t *testing.T) {
	roleTesting := newRoleOrFatal(t, "role::testing")
	roleWWW := newRoleOrFatal(t, "role::www")
	m := NewRoleManager()
	m.Save(roleTesting)
	m.Save(roleWWW)
	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 roles, got %v", len(all))
	}
	if *all[0] != *roleTesting && *all[1] != *roleTesting {
		t.Errorf("missing role: %v", roleTesting)
	}
	if *all[0] != *roleWWW && *all[1] != *roleWWW {
		t.Errorf("missing role: %v", roleWWW)
	}
}

func TestSaveTwiceAndRetrieve(t *testing.T) {
	role := newRoleOrFatal(t, "role::testing")
	m := NewRoleManager()
	m.Save(role)
	m.Save(role)
	all := m.All()
	if len(all) != 1 {
		t.Errorf("expected 1 role, got %v", len(all))
	}
	if *all[0] != *role {
		t.Errorf("expected role %v, got %v", role, all[0])
	}
}

func TestSaveAndFind(t *testing.T) {
	role := newRoleOrFatal(t, "role::testing")
	m := NewRoleManager()
	m.Save(role)
	nt, ok := m.Find(role.Name)
	if !ok {
		t.Errorf("didn't find role")
	}
	if *role != *nt {
		t.Errorf("expected %v, got %v", role, nt)
	}
}
