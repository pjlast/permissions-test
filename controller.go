package main

import (
	"database/sql"
)

type Controller struct {
	DB *sql.DB
}

func (c *Controller) CreateUser(name string) (*User, error) {
	u := User{}
	err := c.DB.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id, name;", name).Scan(&u.ID, &u.Name)

	return &u, err
}

func (c *Controller) CreateRole(name string) (*Role, error) {
	r := Role{}
	err := c.DB.QueryRow("INSERT INTO roles (name) VALUES ($1) RETURNING id, name;", name).Scan(&r.ID, &r.Name)

	return &r, err
}
