package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lunatictiol/that-pet-place-backend-go/services/pets"
	"github.com/lunatictiol/that-pet-place-backend-go/services/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiServer struct {
	address    string
	userdb     *sql.DB
	petStoreDB *mongo.Client
}

func (a *ApiServer) New(address string,
	db *sql.DB, petDb *mongo.Client) {
	a.address = address
	a.userdb = db
	a.petStoreDB = petDb

}

func (a *ApiServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := users.NewStore(a.userdb)
	petStore := pets.NewStore(a.petStoreDB)

	pethandler := pets.NewHandler(petStore)
	userHandler := users.NewHandler(userStore)

	pethandler.RegisterRoutes(subRouter)
	userHandler.RegisterRoutes(subRouter)
	log.Println("listening on port", a.address)
	return http.ListenAndServe(a.address, router)
}
