package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// here we read from the request context and fetch out `"user"` key set in
	// the MyMiddleware example above.
	rawUser := r.Context().Value(userKey)

	if rawUser != nil {
		user := rawUser.(*User)
		fmt.Fprintf(w, "Hello %s!\n", user.Name)
		return
	}

	fmt.Fprintf(w, "Hello World!\n")
}

func getBatchChangesHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(*User)

	rows, err := db.Query("SELECT id, name, private, namespace_org_id, namespace_user_id, creator_id FROM batch_changes WHERE (namespace_user_id = $1) OR (namespace_user_id <> $1 AND private = false) OR (id IN (SELECT namespace_object_id FROM permissions p WHERE p.namespace = 'BATCHCHANGES' AND p.relation = 'VIEW' AND p.namespace_user_id = $1)) OR (EXISTS (SELECT 1 FROM org_members WHERE org_id = batch_changes.namespace_org_id AND user_id = $1 AND org_id <> 0))", user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	bcs := []*batchChange{}
	for rows.Next() {
		var bc batchChange
		err := rows.Scan(&bc.ID, &bc.Name, &bc.Private, &NullInt{N: &bc.NamespaceOrgID}, &NullInt{N: &bc.NamespaceUserID}, &bc.CreatorID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		bcs = append(bcs, &bc)
	}

	render.JSON(w, r, bcs)
	return
}
