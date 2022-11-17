package main

import (
	"fmt"
)

type Permission struct {
	ID        int
	namespace string
	ObjectID  int
	Relation  string
	SubjectID int
}

func (p *Permission) String() string {
	s := fmt.Sprintf("%s:%d#%s", p.namespace, p.ObjectID, p.Relation)

	if p.SubjectID != 0 {
		s += fmt.Sprintf("@%s", p.SubjectID)
	}

	return s
}

func (p *Permission) HasNamespaceAccess(namespace string) bool {
	fmt.Printf("checking namespace access for %s\n", namespace)
	return false
}

//select * from batch_changes bc where (private = true AND ) AND (public = false AND bc.owner_id = 2)
