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
	fmt.Println("============================= SETUP ========================================")
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

	roles, err := seedRoles(c)
	if err != nil {
		log.Fatalf("error while seeding users --> %v", err)
	}
	fmt.Printf("Successfully seeded %d roles ... \n", len(roles))

	s, err := ParseSchema()
	if err != nil {
		log.Fatal("error parsing schema")
	}

	ps, err := c.CreatePermissions(s)
	if err != nil {
		log.Fatalf("error creating permissions ", err)
	}
	for _, r := range roles {
		for _, p := range ps {
			// add all global permissions to the default role and site admin role
			err = c.AddPermissionToRole(r, p)
			if err != nil {
				log.Fatalf("error adding permission '%s' to role %s", p.String(), r.Name)
			}
		}
	}

	fmt.Println("We grab the def")
	defaultRole, err := c.GetRoleByName("DEFAULT")
	if err != nil {
		log.Fatalf("error getting default role ", err)
	}

	users, err := seedUsers(c, defaultRole)
	if err != nil {
		log.Fatalf("error while seeding users --> %v", err)
	}
	fmt.Printf("Successfully seeded %d users ... \n", len(users))
	fmt.Println("============================= SETUP COMPLETE========================================")
}
