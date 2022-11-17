package main

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

const DatabaseURL = "postgres://sourcegraph@localhost:5432/perms?sslmode=disable"

func main() {
	teardown, err := setupDB(DatabaseURL)
	if err != nil {
		log.Fatalf("Migrations: %v", err)
	}
	defer teardown()

	db, err := NewDB(DatabaseURL)
	if err != nil {
		log.Fatal("Getting db connector", err)
	}

	c := &Controller{DB: db}

	users, err := seedUsers(c)
	if err != nil {
		log.Fatalf("error while seeding users --> %v", err)
	}
	fmt.Printf("Successfully seeded %d users ... \n", len(users))

	roles, err := seedRoles(c)
	if err != nil {
		log.Fatalf("error while seeding users --> %v", err)
	}
	fmt.Printf("Successfully seeded %d roles ... \n", len(roles))

}

func seedUsers(c *Controller) ([]*User, error) {
	users := make([]*User, len(mockUserNames))

	for idx, name := range mockUserNames {
		user, err := c.CreateUser(name)
		if err != nil {
			return users, err

		}

		users[idx] = user
	}

	return users, nil
}

func seedRoles(c *Controller) ([]*Role, error) {
	roles := make([]*Role, len(mockRoles))

	for idx, name := range mockRoles {
		role, err := c.CreateRole(name)
		if err != nil {
			return roles, err
		}

		roles[idx] = role
	}

	return roles, nil
}
