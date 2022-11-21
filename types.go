package main

type batchChange struct {
	ID              int
	Name            string
	Private         bool
	NamespaceUserID int
	NamespaceOrgID  int
	CreatorID       int
}

type notebook struct {
	ID        int
	Name      string
	Content   string
	Private   bool
	CreatorID int
}

type codeinsights struct {
	ID     int
	Name   string
	UserID int
}
