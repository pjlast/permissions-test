package main

import "fmt"

// User represents a user of the application.
type User struct {
	ID   int
	Name string
}

var mockUserNames = []string{
	"kai",
	"hunter",
	"elliot",
	"asa",
	"jalen",
	"evan",
	"jude",
	"aubrey",
}

var userRoleMapping = map[string][]string{
	"kai":    {"DEFAULT", "SITE_ADMINISTRATOR"},
	"hunter": {"DEFAULT", "OPERATOR"},
	"elliot": {"DEFAULT", "OPERATOR"},
	"asa":    {"DEFAULT", "SITE_ADMINISTRATOR"},
	"jalen":  {"DEFAULT"},
	"evan":   {"DEFAULT"},
	"jude":   {"DEFAULT"},
	"aubrey": {"DEFAULT"},
}

func seedUsers(c *Controller, roles ...*Role) ([]*User, error) {
	users := make([]*User, len(mockUserNames), len(mockUserNames))

	for idx, name := range mockUserNames {
		user, err := c.createUser(name)
		if err != nil {
			return users, err
		}

		// add users (Kai and Aubrey) to the Acme org
		if name == "kai" || name == "aubrey" {
			fmt.Printf("Adding user %s to Acme org.", name)
			org, err := c.getOrgByName("ACME")
			if err != nil {
				return users, err
			}
			err = c.addUserToOrg(user, org)
			if err != nil {
				return users, err
			}
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
