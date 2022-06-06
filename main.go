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

	m.Up()

	return func() { m.Down() }, nil
}

func createRelation(db *sql.DB, objectNamespace string, objectID int, relation, subjectNamespace string, subjectID int, subjectSetRelation string) error {
	_, err := db.Exec("INSERT INTO permissions(object_namespace, object_id, relation, subject_namespace, subject_id, subject_set_relation) VALUES($1, $2, $3, $4, $5, $6);", objectNamespace, objectID, relation, subjectNamespace, subjectID, subjectSetRelation)

	return err
}

func deleteRelation(db *sql.DB, objectNamespace string, objectID int, relation, subjectNamespace string, subjectID int, subjectSetRelation string) error {
	_, err := db.Exec("DELETE FROM permissions WHERE object_namespace = $1 AND object_id = $2 AND relation = $3 AND subject_namespace = $4 AND subject_id = $5 AND subject_set_relation = $6;", objectNamespace, objectID, relation, subjectNamespace, subjectID, subjectSetRelation)

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
		and relation = $4
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
) select subject_id from all_users where subject_namespace = $3;
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

	fmt.Println("First we show that, without any permission entries, user 1 cannot access notebook 1")

	canRead, err := checkRelation(db, "notebooks", 1, "read", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 read notebook 1?: ", canRead)

	fmt.Println("Now we create an entry giving user 1 direct read access to notebook 1, and check again.")

	err = createRelation(db, "notebooks", 1, "read", "users", 1, "")
	if err != nil {
		panic(err)
	}

	canRead, err = checkRelation(db, "notebooks", 1, "read", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 read notebook 1?: ", canRead)

	fmt.Println("-------")

	fmt.Println("User 1 should still be unable to write to notebook 1.")

	canWrite, err := checkRelation(db, "notebooks", 1, "write", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 write to notebook 1?: ", canWrite)

	fmt.Println("We are now going to give user 1 indirect write access to notebook 1. We will add user 1 to a group, and then give members of that group write access to notebook 1")

	err = createRelation(db, "groups", 1, "member", "users", 1, "")
	err = createRelation(db, "notebooks", 1, "write", "groups", 1, "member")

	fmt.Println("Now, with user 1 part of group 1 which has write access to notebooks 1, we check again:")

	canWrite, err = checkRelation(db, "notebooks", 1, "write", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 write to notebook 1?: ", canWrite)

	fmt.Println("---------")

	fmt.Println("Next we want to show that we can link a code insight's permissions to the notebook, and by extension give user 1 access to the code insight.")

	fmt.Println("First we show that, without any permission entries, user 1 cannot access the code insight")

	canRead, err = checkRelation(db, "codeinsights", 1, "read", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 read codeinsight 1?: ", canRead)

	fmt.Println("Now we give readers of notebook 1 read access to code insight 1, and check again")

	err = createRelation(db, "codeinsights", 1, "read", "notebooks", 1, "read")

	canRead, err = checkRelation(db, "codeinsights", 1, "read", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 read codeinsight 1?: ", canRead)

	fmt.Println("------")

	fmt.Println("Further more, we show that we can give writers of notebook 1 write access to code insight 1.")
	fmt.Println("In this case, user 1 can now write to code insight 1, because he is a member of group 1 which has write access to notebook 1, which now has write access to code insight 1")

	err = createRelation(db, "codeinsights", 1, "write", "notebooks", 1, "write")

	canWrite, err = checkRelation(db, "codeinsights", 1, "write", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 write to codeinsight 1?: ", canWrite)

	fmt.Println("------")
	fmt.Println("Finally, we remove user 1 from group 1, which should cause them to lose write access to notebook 1, but they will still be able to read it and code insight 1")

	err = deleteRelation(db, "groups", 1, "member", "users", 1, "")
	if err != nil {
		panic(err)
	}

	canWrite, err = checkRelation(db, "codeinsights", 1, "write", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 write to codeinsight 1?: ", canWrite)

	canWrite, err = checkRelation(db, "notebooks", 1, "write", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 write to notebook 1?: ", canWrite)

	canRead, err = checkRelation(db, "codeinsights", 1, "read", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 read codeinsight 1?: ", canRead)

	canRead, err = checkRelation(db, "notebooks", 1, "read", "users", 1)

	if err != nil {
		panic(err)
	}

	fmt.Println("Can user 1 read notebook 1?: ", canRead)
}
