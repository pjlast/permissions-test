package main

import "fmt"

// User represents a user of the application.
type User struct {
	ID   int
	Name string
}

func (u *User) checkPublicResourceAccess(namespace, relation string) (bool, error) {
	var id int
	err := db.QueryRow("SELECT id FROM permissions WHERE namespace = $1 AND relation = $2 AND user_id = $3", namespace, relation, u.ID).Scan(&id)
	if err != nil {
		return false, err
	}
	return id > 0, nil
}

var mockUserNames = []string{
	"kai",
	"elliot",
	"jalen",
}

var userRoleMapping = map[string][]string{
	"kai":    {"DEFAULT", "SITE_ADMINISTRATOR"},
	"elliot": {"DEFAULT", "OPERATOR"},
	"jalen":  {"DEFAULT"},
}

func seedUsers(c *Controller, roles ...*Role) ([]*User, error) {
	users := make([]*User, len(mockUserNames), len(mockUserNames))

	for idx, name := range mockUserNames {
		user, err := c.createUser(name)
		if err != nil {
			return users, err
		}

		// add users Kai and Jalen to the Acme org
		if name == "kai" || name == "jalen" {
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
