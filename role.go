package main

import "fmt"

// Role represents a single role in the application.
type Role struct {
	ID          int32
	Name        string
	Permissions []*string
}

var mockRoles = []string{
	"DEFAULT",
	"SITE_ADMINISTRATOR",
	"OPERATOR",
}

func seedRoles(c *Controller) ([]*Role, error) {
	fmt.Println("")
	roles := make([]*Role, len(mockRoles), len(mockRoles))

	for idx, name := range mockRoles {
		fmt.Printf("Seeding role '%s'.\n", name)
		role, err := c.createRole(name)
		if err != nil {
			return roles, err
		}

		roles[idx] = role
	}

	return roles, nil
}
