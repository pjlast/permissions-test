package main

import (
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const databaseURL = "postgres://sourcegraph@localhost:5432/perms?sslmode=disable"

func main() {
	fmt.Println("============================= SETUP ========================================")
	teardown, err := setupDB(databaseURL)
	if err != nil {
		log.Fatalf("Migrations: %v", err)
	}
	defer teardown()

	db, err := newDB(databaseURL)
	if err != nil {
		log.Fatal("Getting db connector", err)
	}

	c := &Controller{DB: db}

	rs, err := seedRoles(c)
	if err != nil {
		log.Fatalf("error while seeding users --> %s", err.Error())
	}
	fmt.Printf("Successfully seeded %d roles ... \n", len(rs))

	s, err := ParseSchema()
	if err != nil {
		log.Fatal("error parsing schema")
	}

	ps, err := c.createPermissions(s)
	if err != nil {
		log.Fatalf("error creating permissions: %s", err.Error())
	}

	err = populateRolesWithPermissions(c, rs, ps)
	if err != nil {
		log.Fatalf("error populating roles with permissions: %s", err.Error())
	}

	users, err := seedUsers(c, rs...)
	if err != nil {
		log.Fatalf("error while seeding users --> %v", err)
	}
	fmt.Printf("Successfully seeded %d users ... \n", len(users))
	fmt.Println("============================= SETUP COMPLETE========================================")
}

func populateRolesWithPermissions(c *Controller, rs []*Role, ps []*Permission) error {
	fmt.Println("")
	fmt.Println(`We want the site administrator to have all permissions, however we want
the DEFAULT role to give full access to only notebooks,
view access to batch changes and view to code insights,
while the operator role will give full access to batch changes`)
	fmt.Println("")
	for _, r := range rs {
		// add all global permissions to the site admin role
		if r.Name == "SITE_ADMINISTRATOR" {
			fmt.Println("Adding all global permissions to the site admin role")
			for _, p := range ps {
				err := c.addPermissionToRole(r, p)
				if err != nil {
					return err
				}
			}
			continue
		}

		if r.Name == "DEFAULT" {
			fmt.Println("Adding permissions to the DEFAULT role ...")
			for _, p := range ps {
				if p.Namespace == "NOTEBOOKS" {
					err := c.addPermissionToRole(r, p)
					if err != nil {
						return err
					}
				}

				if p.Namespace == "CODEINSIGHTS" && p.Relation == "VIEW" {
					err := c.addPermissionToRole(r, p)
					if err != nil {
						return err
					}
				}

				if p.Namespace == "BATCHCHANGES" && p.Relation == "VIEW" {
					err := c.addPermissionToRole(r, p)
					if err != nil {
						return err
					}
				}
			}

			continue
		}

		if r.Name == "OPERATOR" {
			fmt.Println("Adding permissions to the OPERATOR role ...")
			for _, p := range ps {
				if p.Namespace == "BATCHCHANGES" {
					err := c.addPermissionToRole(r, p)
					if err != nil {
						return err
					}
				}
			}

			continue
		}
	}

	return nil
}
