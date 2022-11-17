package main

type User struct {
	ID   int
	Name string
}

var mockUserNames = []string{
	"Kai",
	"Hunter",
	"Elliot",
	"Asa",
	"Jalen",
	"Evan",
	"Jude",
	"Aubrey",
}

//func seedUsers(db *sql.DB) error {
//	users := make([]*User, len(mockUserNames))
//	for _, name := range mockUserNames {
//		users = append(users, &User{
//			Name: name,
//		})
//	}
//
//	_, err := db.Exec()
//
//	return nil
//}
//
//var insertUserFmtQuery = `INSERT INTO users
//VALUES (%s)
//RETURNING %s`
