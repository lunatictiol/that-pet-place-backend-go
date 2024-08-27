package main

import (
	"fmt"
	"log"

	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lunatictiol/that-pet-place-backend-go/config"
	"github.com/lunatictiol/that-pet-place-backend-go/db"
)

func main() {
	connStr := fmt.Sprintf("user='%s' password=%s host=%s dbname='%s'", config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBHost, config.Envs.DBName)

	db, err := db.NewMySqlStorage(connStr)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

}
