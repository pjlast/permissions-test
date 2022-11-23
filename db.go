package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func setupDB(url string) (func(), error) {
	m, err := migrate.New(
		"file://migrations",
		url,
	)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully executed the migrations.")

	return func() {
		fmt.Println("")
		fmt.Println("===================== TEARDOWN =====================")
		fmt.Println("Reversing migrations ....")
		err := m.Down()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("===================== TEARDOWN COMPLETE =====================")
	}, nil
}

func newDB(url string) (*sql.DB, error) {
	return sql.Open("postgres", url)
}

// NullInt represents an int that may be null. NullInt implements the
// sql.Scanner interface so it can be used as a scan destination, similar to
// sql.NullString. When the scanned value is null, int is set to the zero value.
type NullInt struct{ N *int }

// NewNullInt returns a NullInt treating zero value as null.
func NewNullInt(i int) NullInt {
	if i == 0 {
		return NullInt{}
	}
	return NullInt{N: &i}
}

// Scan implements the Scanner interface.
func (n *NullInt) Scan(value any) error {
	switch value := value.(type) {
	case int64:
		*n.N = int(value)
	case int32:
		*n.N = int(value)
	case nil:
		return nil
	default:
		return fmt.Errorf("value is not int: %T", value)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullInt) Value() (driver.Value, error) {
	if n.N == nil {
		return nil, nil
	}
	return *n.N, nil
}
