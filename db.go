package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"time"
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
		time.Sleep(10 * time.Second)
		fmt.Println("Reversing migrations ....")
		err := m.Down()
		if err != nil {
			log.Fatal(err)
		}
	}, nil
}

func NewDB(url string) (*sql.DB, error) {
	return sql.Open("postgres", url)
}
