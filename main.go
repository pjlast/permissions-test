package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func setupDB() (func(), error) {
	godotenv.Load()
	m, err := migrate.New(
		"file://db/migration",
		os.Getenv("TEST_DATABASE_URL"),
	)

	if err != nil {
		return nil, err
	}

	m.Up()

	return func() { m.Down() }, nil
}

type permissions struct {
	DB               *sql.DB
	Namespaces       []string
	Relations        map[string]int
	ReverseRelations map[int]string
}

func (p *permissions) createResource(resourceType, resourceName string) (int, error) {
	var id int
	err := p.DB.QueryRow(fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id;", resourceType), resourceName).Scan(&id)

	return id, err
}

func (p *permissions) createUserSet(namespace string, id int, relation string) (int, error) {
	var usersetID int
	err := p.DB.QueryRow(fmt.Sprintf(`
WITH e AS(
	INSERT INTO usersets (relation, %[1]s_id)
		VALUES ($1, $2)
	ON CONFLICT(relation, %[1]s_id) DO NOTHING
	RETURNING id
)
SELECT * FROM e
UNION
SELECT id FROM usersets WHERE relation=$1 AND %[1]s_id=$2`, namespace), p.Relations[relation], id).Scan(&usersetID)

	return usersetID, err
}

func (p *permissions) createUserRelation(namespace string, id int, relation string, userID int) error {
	_, err := p.DB.Exec(fmt.Sprintf("INSERT INTO %s_namespace(id, relation, user_id) VALUES($1, $2, $3);", namespace), id, p.Relations[relation], userID)

	return err
}

func (p *permissions) createUserSetRelation(namespace string, id int, relation string, usersetNamespace string, usersetID int, usersetRelation string) error {
	usersetID, err := p.createUserSet(usersetNamespace, usersetID, usersetRelation)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec(fmt.Sprintf("INSERT INTO %s_namespace(id, relation, userset_id) VALUES($1, $2, $3);", namespace), id, p.Relations[relation], usersetID)

	return err
}

func (p *permissions) deleteUserRelation(namespace string, id int, relation string, userID int) error {
	_, err := p.DB.Exec(fmt.Sprintf("DELETE FROM %s_namespace WHERE id = $1 AND relation = $2 AND user_id = $3;", namespace), id, p.Relations[relation], userID)

	return err
}

func (p *permissions) deleteUserSetRelation(namespace string, id int, relation string, userID int) error {
	_, err := p.DB.Exec(fmt.Sprintf("DELETE FROM %s_namespace WHERE id = $1 AND relation = $2 AND user_id = $3;", namespace), id, p.Relations[relation], userID)

	return err
}

func (p *permissions) checkUserRelation(namespace string, id int, relation string, userID int) (bool, error) {
	idSlice := make([]string, 0, len(p.Namespaces))
	relationSlice := make([]string, 0, len(p.Namespaces))
	userIDSlice := make([]string, 0, len(p.Namespaces))
	usersetIDSlice := make([]string, 0, len(p.Namespaces))
	namespaceJoinsSlice := make([]string, 0, len(p.Namespaces))
	for _, k := range p.Namespaces {
		idSlice = append(idSlice, fmt.Sprintf("%s_namespace.id", k))
		relationSlice = append(relationSlice, fmt.Sprintf("%s_namespace.relation", k))
		userIDSlice = append(userIDSlice, fmt.Sprintf("%s_namespace.user_id", k))
		usersetIDSlice = append(usersetIDSlice, fmt.Sprintf("%s_namespace.userset_id", k))
		namespaceJoinsSlice = append(namespaceJoinsSlice, fmt.Sprintf(`LEFT JOIN %[1]s_namespace ON
			%[1]s_namespace.id = usersets.%[1]s_id AND
			%[1]s_namespace.relation = usersets.relation`, k))
	}

	sqlQuery := fmt.Sprintf(`
WITH RECURSIVE all_users AS (
	SELECT
		id,
		relation,
		user_id,
		userset_id,
		relation AS original_relation
	FROM
		%[1]s_namespace
	WHERE
		id = $1 AND
		relation = $2
	UNION
	SELECT
		COALESCE(%[2]s) as id,
		COALESCE(%[3]s) as relation,
		COALESCE(%[4]s) as user_id,
		COALESCE(%[5]s) as userset_id,
		au.original_relation
	FROM
		all_users au
	INNER JOIN usersets ON
		au.userset_id = usersets.id
	%[6]s
)
SELECT
	user_id
FROM all_users
WHERE
	user_id IS NOT NULL;
`, namespace, strings.Join(idSlice, ", "), strings.Join(relationSlice, ", "), strings.Join(userIDSlice, ", "), strings.Join(usersetIDSlice, ", "), strings.Join(namespaceJoinsSlice, "\n"))

	rows, err := p.DB.Query(sqlQuery, id, p.Relations[relation])

	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		var rowUserID int
		err := rows.Scan(&rowUserID)

		if err != nil {
			return false, err
		}

		if rowUserID == userID {
			return true, nil
		}
	}

	return false, nil
}

func (p *permissions) fetchUsersWithRelation(namespace string, id int, relation string) ([]int, error) {
	idSlice := make([]string, 0, len(p.Namespaces))
	relationSlice := make([]string, 0, len(p.Namespaces))
	userIDSlice := make([]string, 0, len(p.Namespaces))
	usersetIDSlice := make([]string, 0, len(p.Namespaces))
	namespaceJoinsSlice := make([]string, 0, len(p.Namespaces))
	for _, k := range p.Namespaces {
		idSlice = append(idSlice, fmt.Sprintf("%s_namespace.id", k))
		relationSlice = append(relationSlice, fmt.Sprintf("%s_namespace.relation", k))
		userIDSlice = append(userIDSlice, fmt.Sprintf("%s_namespace.user_id", k))
		usersetIDSlice = append(usersetIDSlice, fmt.Sprintf("%s_namespace.userset_id", k))
		namespaceJoinsSlice = append(namespaceJoinsSlice, fmt.Sprintf(`LEFT JOIN %[1]s_namespace ON
			%[1]s_namespace.id = usersets.%[1]s_id AND
			%[1]s_namespace.relation = usersets.relation`, k))
	}

	sqlQuery := fmt.Sprintf(`
WITH RECURSIVE all_users AS (
	SELECT
		id,
		relation,
		user_id,
		userset_id,
		relation AS original_relation
	FROM
		%[1]s_namespace
	WHERE
		id = $1 AND
		relation = $2
	UNION
	SELECT
		COALESCE(%[2]s) as id,
		COALESCE(%[3]s) as relation,
		COALESCE(%[4]s) as user_id,
		COALESCE(%[5]s) as userset_id,
		au.original_relation
	FROM
		all_users au
	INNER JOIN usersets ON
		au.userset_id = usersets.id
	%[6]s
)
SELECT
	user_id
FROM all_users
WHERE
	user_id IS NOT NULL;
`, namespace, strings.Join(idSlice, ", "), strings.Join(relationSlice, ", "), strings.Join(userIDSlice, ", "), strings.Join(usersetIDSlice, ", "), strings.Join(namespaceJoinsSlice, "\n"))

	rows, err := p.DB.Query(sqlQuery, id, p.Relations[relation])

	if err != nil {
		return []int{}, err
	}

	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var rowUserID int
		err := rows.Scan(&rowUserID)

		if err != nil {
			return []int{}, err
		}

		userIDs = append(userIDs, rowUserID)
	}

	return userIDs, nil
}

type permsStruct struct {
	CodeinsightID *int
	NotebookID    *int
	GroupID       *int
	Relation      *int
}

func (p *permissions) fetchResourcesForUser(userID int) ([]permsStruct, error) {
	generateSelectQuery := func(namespaces []string, namespace string) string {
		output := "SELECT\n"
		for _, k := range p.Namespaces {
			if k == namespace {
				output += fmt.Sprintf("id AS %s_id,\n", k)
			} else {
				output += fmt.Sprintf("CAST(NULL AS int) AS %s_id,\n", k)
			}
		}
		output += fmt.Sprintf(`
		relation,
		user_id,
		userset_id
	FROM
		%[1]s_namespace
	WHERE
		%[1]s_namespace.user_id = $1
`, namespace)
		return output
	}

	idSlice := make([]string, 0, len(p.Namespaces))
	relationSlice := make([]string, 0, len(p.Namespaces))
	namespaceSelectsSlice := make([]string, 0, len(p.Namespaces))
	namespaceIDSlice := make([]string, 0, len(p.Namespaces))
	innerjoinIDSlice := make([]string, 0, len(p.Namespaces))
	leftjoinIDSlice := make([]string, 0, len(p.Namespaces))
	for _, k := range p.Namespaces {
		idSlice = append(idSlice, fmt.Sprintf("%s_id", k))
		relationSlice = append(relationSlice, fmt.Sprintf("%s_namespace.relation", k))
		namespaceSelectsSlice = append(namespaceSelectsSlice, generateSelectQuery(p.Namespaces, k))
		namespaceIDSlice = append(namespaceIDSlice, fmt.Sprintf("%[1]s_namespace.id AS %[1]s_id,", k))
		innerjoinIDSlice = append(innerjoinIDSlice, fmt.Sprintf("ar.%[1]s_id = usersets.%[1]s_id", k))
		leftjoinIDSlice = append(leftjoinIDSlice, fmt.Sprintf(`LEFT JOIN %[1]s_namespace ON
		%[1]s_namespace.userset_id = usersets.id`, k))
	}

	sqlQuery := fmt.Sprintf(`
	WITH RECURSIVE all_resources AS (
		(%[1]s
		)
		UNION
		select
			%[2]s
			coalesce(%[3]s) as relation,
			null as user_id,
			null as userset_id
		FROM
			all_resources ar
		INNER JOIN usersets ON
			(%[4]s) AND ar.relation = usersets.relation
		%[5]s
	)
	SELECT
		%[6]s,
		relation
	FROM all_resources;
`, strings.Join(namespaceSelectsSlice, "UNION\n"), strings.Join(namespaceIDSlice, "\n"), strings.Join(relationSlice, ", "), strings.Join(innerjoinIDSlice, " OR\n"), strings.Join(leftjoinIDSlice, "\n"), strings.Join(idSlice, ",\n"))

	rows, err := p.DB.Query(sqlQuery, userID)

	if err != nil {
		return []permsStruct{}, err
	}

	defer rows.Close()

	var rowPermissions []permsStruct

	for rows.Next() {
		var rowPerm permsStruct
		err := rows.Scan(&rowPerm.NotebookID, &rowPerm.CodeinsightID, &rowPerm.GroupID, &rowPerm.Relation)

		if err != nil {
			return []permsStruct{}, err
		}

		rowPermissions = append(rowPermissions, rowPerm)
	}

	return rowPermissions, nil
}

func main() {
	teardown, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer teardown()

	db, err := sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	perms := permissions{DB: db, Relations: make(map[string]int), ReverseRelations: make(map[int]string)}

	// Create some relations we can reference
	id, err := perms.createResource("relations", "view")
	if err != nil {
		panic(err)
	}
	perms.Relations["view"] = id
	perms.ReverseRelations[id] = "view"

	id, err = perms.createResource("relations", "edit")
	if err != nil {
		panic(err)
	}
	perms.Relations["edit"] = id
	perms.ReverseRelations[id] = "edit"

	// Load the namespaces
	perms.Namespaces = []string{"notebooks", "codeinsights", "groups"}

	// Create a user
	steveID, _ := perms.createResource("users", "Steve")
	// Create a notebook
	notebook1ID, _ := perms.createResource("notebooks", "Notebook 1")

	fmt.Println("First we show that, without any permission entries, user Steve cannot view notebook 1")

	canView, err := perms.checkUserRelation("notebooks", notebook1ID, "view", steveID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Can user Steve view notebook 1?: ", canView)

	fmt.Println("Now we create an entry giving user Steve direct view access to notebook 1, and check again.")

	err = perms.createUserRelation("notebooks", notebook1ID, "view", steveID)

	if err != nil {
		panic(err)
	}

	canView, err = perms.checkUserRelation("notebooks", notebook1ID, "view", steveID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Can user Steve view notebook 1?: ", canView)

	fmt.Println("-------")

	fmt.Println("User Steve should still be unable to edit to notebook 1.")

	canEdit, err := perms.checkUserRelation("notebooks", notebook1ID, "edit", steveID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Can user Steve edit notebook 1?: ", canEdit)

	fmt.Println("We are now going to give user Steve indirect edit access to notebook 1. We will add user Steve to a group, and then give members of that group edit access to notebook 1")

	fmt.Println("First we create a group, define the 'member' relation, and create a member relation between Steve and the group")
	groupID, err := perms.createResource("groups", "Group 1")
	if err != nil {
		panic(err)
	}

	relationID, err := perms.createResource("relations", "member")
	if err != nil {
		panic(err)
	}
	perms.Relations["member"] = relationID
	perms.ReverseRelations[relationID] = "member"

	err = perms.createUserRelation("groups", groupID, "member", steveID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Now we create an edit relation to the userset that is 'members of group 1'")
	err = perms.createUserSetRelation("notebooks", notebook1ID, "edit", "groups", groupID, "member")
	if err != nil {
		panic(err)
	}

	fmt.Println("Now, with user Steve part of group 1 which has edit access to notebook 1, we check again:")

	canEdit, err = perms.checkUserRelation("notebooks", notebook1ID, "edit", steveID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Can user Steve edit notebook 1?: ", canEdit)

	fmt.Println("---------")

	fmt.Println("Next we want to show that we can link a code insight's permissions to the notebook, and by extension give user Steve access to the code insight.")
	fmt.Println("First we create a code insight")

	codeinsightID, err := perms.createResource("codeinsights", "Code Insight 1")
	if err != nil {
		panic(err)
	}

	fmt.Println("Now we show that, without any permission entries, user 1 cannot access the code insight")

	canView, err = perms.checkUserRelation("codeinsights", codeinsightID, "view", steveID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Can user Steve view Code Insight 1?: ", canView)

	fmt.Println("Now we give viewers of notebook 1 view access to code insight 1, and check again")

	err = perms.createUserSetRelation("codeinsights", codeinsightID, "view", "notebooks", notebook1ID, "view")
	if err != nil {
		panic(err)
	}

	canView, err = perms.checkUserRelation("codeinsights", codeinsightID, "view", steveID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Can user Steve view Code Insight 1?: ", canView)

	fmt.Println("------")

	fmt.Println("Further more, we show that we can give editors of notebook 1 edit access to code insight 1.")
	fmt.Println("In this case, user Steve can now edit code insight 1, because he is a member of group 1 which has edit access to notebook 1, which now has edit access to code insight 1")

	err = perms.createUserSetRelation("codeinsights", codeinsightID, "edit", "notebooks", notebook1ID, "edit")

	canEdit, err = perms.checkUserRelation("codeinsights", codeinsightID, "edit", steveID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Can user Steve edit Code Insight 1?: ", canEdit)

	fmt.Println("------")

	fmt.Println("We also want to show that we are able to retrieve all users that have access to a resource, or all resources that a user has access to.")
	fmt.Println("First we will give another user, user Sarah, direct write access to the code insight.")

	sarahID, err := perms.createResource("users", "Sarah")
	if err != nil {
		panic(err)
	}

	err = perms.createUserRelation("codeinsights", codeinsightID, "edit", sarahID)

	fmt.Println("Then we will query for all users that have edit access to code insight 1:")

	userIDs, err := perms.fetchUsersWithRelation("codeinsights", codeinsightID, "edit")

	fmt.Println("User IDs with write access to code insight 1 (1 is Steve, 2 is Sarah): ", userIDs)

	fmt.Println("-----")
	fmt.Println("Similarly, we can fetch all resources a user has access to. We'll create a code insight 2 that user Steve can also edit, and retrieve the list")

	codeInsight2ID, err := perms.createResource("codeinsights", "Code Insight 2")
	if err != nil {
		panic(err)
	}

	err = perms.createUserRelation("codeinsights", codeInsight2ID, "edit", steveID)
	if err != nil {
		panic(err)
	}

	stevePerms, err := perms.fetchResourcesForUser(steveID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Permissions for user Steve: ")
	for _, perm := range stevePerms {
		if perm.CodeinsightID != nil {
			fmt.Println(fmt.Sprintf("Codeinsight %d - %s", *perm.CodeinsightID, perms.ReverseRelations[*perm.Relation]))
		}
		if perm.NotebookID != nil {
			fmt.Println(fmt.Sprintf("Notebook %d - %s", *perm.NotebookID, perms.ReverseRelations[*perm.Relation]))
		}
		if perm.GroupID != nil {
			fmt.Println(fmt.Sprintf("Group %d - %s", *perm.GroupID, perms.ReverseRelations[*perm.Relation]))
		}
	}

	fmt.Println("-----")
	fmt.Println("Finally, we remove user 1 from group 1, which should cause them to lose write access to notebook 1, but they will still be able to read it and code insight 1")

	err = perms.deleteUserRelation("groups", groupID, "member", steveID)
	if err != nil {
		panic(err)
	}

	stevePerms, err = perms.fetchResourcesForUser(steveID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Permissions for user Steve: ")
	for _, perm := range stevePerms {
		if perm.CodeinsightID != nil {
			fmt.Println(fmt.Sprintf("Codeinsight %d - %s", *perm.CodeinsightID, perms.ReverseRelations[*perm.Relation]))
		}
		if perm.NotebookID != nil {
			fmt.Println(fmt.Sprintf("Notebook %d - %s", *perm.NotebookID, perms.ReverseRelations[*perm.Relation]))
		}
		if perm.GroupID != nil {
			fmt.Println(fmt.Sprintf("Group %d - %s", *perm.GroupID, perms.ReverseRelations[*perm.Relation]))
		}
	}
}
