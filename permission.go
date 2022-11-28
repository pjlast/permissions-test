package main

import (
	"fmt"
)

type Permission struct {
	ID                int
	Namespace         string
	NamespaceObjectID int
	Action            string
	NamespaceUserID   int
	NamespaceOrgID    int
}

func (p *Permission) String() string {
	if p.NamespaceUserID != 0 {
		return fmt.Sprintf("%s:%d#%s@%d", p.Namespace, p.NamespaceObjectID, p.Action, p.NamespaceUserID)
	}

	if p.NamespaceOrgID == 0 {
		return fmt.Sprintf("%s:%d#%s@%d", p.Namespace, p.NamespaceObjectID, p.Action, p.NamespaceOrgID)
	}

	return fmt.Sprintf("%s:*#%s", p.Namespace, p.Action)
}

// select * from batch_changes bc where (private = true AND ) AND (public = false AND bc.owner_id = 2)
