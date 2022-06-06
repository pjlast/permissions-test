package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	m.Steps(1)

	return func() { m.Steps(-1) }, nil
}

func createRelation(db *sql.DB, objectNamespace string, objectID int, relation, subjectNamespace string, subjectID int, subjectSetRelation string) error {
	_, err := db.Exec("INSERT INTO permissions(object_namespace, object_id, relation, subject_namespace, subject_id, subject_set_relation) VALUES($1, $2, $3, $4, $5, $6);", objectNamespace, objectID, relation, subjectNamespace, subjectID, subjectSetRelation)

	return err
}

func checkRelation(db *sql.DB, objectNamespace string, objectID int, relation, subjectNamespace string, subjectID int) (bool, error) {
	sqlQuery := `
with recursive all_users as (
	select 
		object_namespace,
		object_id,
		relation,
		subject_namespace,
		subject_id,
		subject_set_relation
	from
		permissions
	where 
		object_namespace = $1
		and object_id = $2
	union 
		select 
			p.object_namespace,
			p.object_id,
			p.relation,
			p.subject_namespace,
			p.subject_id,
			p.subject_set_relation
		from
			permissions p
		inner join all_users au on au.subject_namespace = p.object_namespace and au.subject_id = p.object_id and au.subject_set_relation=p.relation
) select subject_id from all_users where subject_namespace = $3 AND relation = $4;
`

	rows, err := db.Query(sqlQuery, objectNamespace, objectID, subjectNamespace, relation)

	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)

		if err != nil {
			return false, err
		}

		if id == subjectID {
			return true, nil
		}
	}

	return false, nil
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

	err = createRelation(db, "notebooks", 1, "read", "users", 1, "")
	if err != nil {
		panic(err)
	}

	canRead, err := checkRelation(db, "notebooks", 1, "read", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 read notebook 1?: ", canRead)
}
