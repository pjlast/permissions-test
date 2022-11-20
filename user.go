package main

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

func seedUsers(c *Controller, defaultRole *Role) ([]*User, error) {
	users := make([]*User, len(mockUserNames), len(mockUserNames))

	for idx, name := range mockUserNames {
		user, err := c.CreateUser(name)
		if err != nil {
			return users, err
		}

		err = c.AddRoleForUser(user, defaultRole)
		if err != nil {
			return users, err
		}

		users[idx] = user
	}

	return users, nil
}
