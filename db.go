package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
		time.Sleep(30 * time.Second)
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
