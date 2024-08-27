package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lunatictiol/that-pet-place-backend-go/services/users"
)

type ApiServer struct {
	address string
	db      *sql.DB
}

func (a *ApiServer) New(address string,
	db *sql.DB) {
	a.address = address
	a.db = db

}

func (a *ApiServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	userStore := users.NewStore(a.db)
	userHandler := users.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)
	log.Println("listening on port", a.address)
	return http.ListenAndServe(a.address, router)
}
