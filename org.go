package main

import "fmt"

// Org represents an organization
type Org struct {
	ID   int
	Name string
}

var mockOrgs = []string{
	"ACME",
	"HORSEGRAPH",
}

func seedOrgs(c *Controller) ([]*Org, error) {
	fmt.Println("")
	orgs := make([]*Org, len(mockOrgs), len(mockOrgs))

	for idx, name := range mockOrgs {
		fmt.Printf("Seeding org '%s'.\n", name)
		org, err := c.createOrg(name)
		if err != nil {
			return orgs, err
		}

		orgs[idx] = org
	}

	return orgs, nil
}
