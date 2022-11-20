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

func (c *Controller) createPermissions(s *Schema) (ps []*Permission, err error) {
	fmt.Println("")
	for _, namespace := range s.Namespaces {
		fmt.Printf("creating permissions for namespace '%s'.\n", namespace.Name)

		for _, relation := range namespace.Relations {
			p := &Permission{
				Relation:  relation,
				Namespace: namespace.Name,
			}
			err := c.DB.QueryRow("INSERT INTO permissions (namespace, relation) VALUES ($1, $2) RETURNING id;", namespace.Name, relation).Scan(&p.ID)
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
