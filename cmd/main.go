package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lunatictiol/that-pet-place-backend-go/cmd/api"
	"github.com/lunatictiol/that-pet-place-backend-go/config"
	"github.com/lunatictiol/that-pet-place-backend-go/db"
)

func main() {
	var apiServer api.ApiServer
	connStr := fmt.Sprintf("user='%s' password=%s host=%s dbname='%s'", config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBHost, config.Envs.DBName)

	db, err := db.NewMySqlStorage(connStr)
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)
	apiServer.New(":8000", db)
	err = apiServer.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func initStorage(d *sql.DB) {
	err := d.Ping()
	log.Println("database connection succesfull")
	if err != nil {
		log.Fatal(err)
	}

}
