package main

type batchChange struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Private         bool   `json:"private"`
	NamespaceUserID int    `json:"namespace_user_id"`
	NamespaceOrgID  int    `json:"namespace_org_id"`
	CreatorID       int    `json:"creator_id"`
}

type notebook struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Private   bool   `json:"private"`
	CreatorID int    `json:"creator_id"`
}

type codeinsight struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}
