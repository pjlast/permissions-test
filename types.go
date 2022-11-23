package main

type batchChange struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Private         bool   `json:"private"`
	NamespaceUserID int    `json:"namespace_user_id"`
	NamespaceOrgID  int    `json:"namespace_org_id"`
	CreatorID       int    `json:"creator_id"`
}

func (b *batchChange) shareResourceAccess(recipientUserID int, relation string) error {
	var pint int
	err := db.QueryRow("INSERT INTO permissions (namespace, namespace_object_id, namespace_user_id, relation) VALUES ('BATCHCHANGES', $1, $2, $3) RETURNING id", b.ID, recipientUserID, relation).Scan(&pint)

	return err
}

type notebook struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Private   bool   `json:"private"`
	CreatorID int    `json:"creator_id"`
}

func (n *notebook) shareResourceAccess(recipientUserID int, relation string) error {
	var pint int
	err := db.QueryRow("INSERT INTO permissions (namespace, namespace_object_id, namespace_user_id, relation) VALUES ('BATCHCHANGES', $1, $2, $3) RETURNING id", n.ID, recipientUserID, relation).Scan(&pint)

	return err
}

type codeinsight struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}

func (c *codeinsight) shareResourceAccess(recipientUserID int, relation string) error {
	var pint int
	err := db.QueryRow("INSERT INTO permissions (namespace, namespace_object_id, namespace_user_id, relation) VALUES ('BATCHCHANGES', $1, $2, $3) RETURNING id", c.ID, recipientUserID, relation).Scan(&pint)

	return err
}
