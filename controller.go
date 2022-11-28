package main

import (
	"database/sql"
	"fmt"
)

// Controller is the main controller for the application.
type Controller struct {
	DB *sql.DB
}

func (c *Controller) createUser(name string) (*User, error) {
	u := User{}
	err := c.DB.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id, name;", name).Scan(&u.ID, &u.Name)

	return &u, err
}

func (c *Controller) createRole(name string) (*Role, error) {
	r := Role{}
	err := c.DB.QueryRow("INSERT INTO roles (name) VALUES ($1) RETURNING id, name;", name).Scan(&r.ID, &r.Name)

	return &r, err
}

func (c *Controller) createOrg(name string) (*Org, error) {
	o := Org{}
	err := c.DB.QueryRow("INSERT INTO orgs (name) VALUES ($1) RETURNING id, name;", name).Scan(&o.ID, &o.Name)

	return &o, err
}

func (c *Controller) createPermissions(s *Schema) (ps []*Permission, err error) {
	fmt.Println("")
	for _, namespace := range s.Namespaces {
		fmt.Printf("creating permissions for namespace '%s'.\n", namespace.Name)

		for _, action := range namespace.Actions {
			p := &Permission{
				Action:    action,
				Namespace: namespace.Name,
			}
			err := c.DB.QueryRow("INSERT INTO permissions (namespace, action) VALUES ($1, $2) RETURNING id;", namespace.Name, action).Scan(&p.ID)
			if err != nil {
				return ps, err
			}
			ps = append(ps, p)
		}
	}
	return ps, err
}

func (c *Controller) addPermissionToRole(r *Role, p *Permission) error {
	row := c.DB.QueryRow("INSERT INTO role_permissions (permission_id, role_id) VALUES ($1, $2);", p.ID, r.ID)
	return row.Err()
}

func (c *Controller) getRoleByName(name string) (_ *Role, err error) {
	var r = Role{}
	err = c.DB.QueryRow("SELECT id, name FROM roles WHERE name = $1;", name).Scan(&r.ID, &r.Name)
	return &r, err
}

func (c *Controller) addRoleForUser(u *User, rs ...*Role) error {
	for _, r := range rs {
		fmt.Printf("User '%s' has the '%s' role.\n", u.Name, r.Name)
		row := c.DB.QueryRow("INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2);", u.ID, r.ID)
		err := row.Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) addUserToOrg(u *User, o *Org) error {
	row := c.DB.QueryRow("INSERT INTO org_members (user_id, org_id) VALUES ($1, $2);", u.ID, o.ID)
	return row.Err()
}

func (c *Controller) getOrgByName(orgName string) (*Org, error) {
	o := Org{}
	err := c.DB.QueryRow("SELECT id, name FROM orgs WHERE name = $1;", orgName).Scan(&o.ID, &o.Name)

	return &o, err
}

func (c *Controller) createBatchChange(bc *batchChange) error {
	if bc.NamespaceOrgID == 0 && bc.NamespaceUserID == 0 {
		return fmt.Errorf("batch change must have a namespace")
	}

	if bc.NamespaceOrgID != 0 {
		return c.DB.QueryRow("INSERT INTO batch_changes (name, private, namespace_org_id, creator_id) VALUES ($1, $2, $3, $4) RETURNING id;", bc.Name, bc.Private, bc.NamespaceOrgID, bc.CreatorID).Scan(&bc.ID)
	}

	return c.DB.QueryRow("INSERT INTO batch_changes (name, private, namespace_user_id, creator_id) VALUES ($1, $2, $3, $4) RETURNING id;", bc.Name, bc.Private, bc.NamespaceUserID, bc.CreatorID).Scan(&bc.ID)
}

func (c *Controller) createNotebook(n *notebook) error {
	err := c.DB.QueryRow("INSERT INTO notebooks (name, content, private, creator_id) VALUES ($1, $2, $3, $4) RETURNING id;", n.Name, n.Content, n.Private, n.CreatorID).Scan(&n.ID)

	return err
}

func (c *Controller) createCodeInsight(ci *codeinsight) error {
	err := c.DB.QueryRow("INSERT INTO code_insights (name, user_id) VALUES ($1, $2) RETURNING id;", ci.Name, ci.UserID).Scan(&ci.ID)

	return err
}
