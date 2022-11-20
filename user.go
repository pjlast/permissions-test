package main

import "fmt"

// User represents a user of the application.
type User struct {
	ID   int
	Name string
}

var mockUserNames = []string{
	"Kai",
	"Hunter",
	"Elliot",
	"Asa",
	"Jalen",
	"Evan",
	"Jude",
	"Aubrey",
}

var userRoleMapping = map[string][]string{
	"Kai":    {"DEFAULT", "SITE_ADMINISTRATOR"},
	"Hunter": {"DEFAULT", "OPERATOR"},
	"Elliot": {"DEFAULT", "OPERATOR"},
	"Asa":    {"DEFAULT", "SITE_ADMINISTRATOR"},
	"Jalen":  {"DEFAULT"},
	"Evan":   {"DEFAULT"},
	"Jude":   {"DEFAULT"},
	"Aubrey": {"DEFAULT"},
}

func seedUsers(c *Controller, roles ...*Role) ([]*User, error) {
	users := make([]*User, len(mockUserNames), len(mockUserNames))

	for idx, name := range mockUserNames {
		user, err := c.createUser(name)
		if err != nil {
			return users, err
		}

		fmt.Println("")
		roleNames := userRoleMapping[name]
		roles := getUserRoles(roleNames, roles)

		err = c.addRoleForUser(user, roles...)
		if err != nil {
			return users, err
		}

		users[idx] = user
	}

	return users, nil
}

func getUserRoles(roleNames []string, roles []*Role) (rs []*Role) {
	for _, r := range roles {
		for _, rn := range roleNames {
			if r.Name == rn {
				rs = append(rs, r)
			}
		}
	}

	return
}
