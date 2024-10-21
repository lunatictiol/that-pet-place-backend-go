package api

import (
	"database/sql"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	localstores "github.com/lunatictiol/that-pet-place-backend-go/services/localStores"
	"github.com/lunatictiol/that-pet-place-backend-go/services/pets"
	"github.com/lunatictiol/that-pet-place-backend-go/services/users"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiServer struct {
	address     string
	userdb      *sql.DB
	petsDB      *sql.DB
	shopeStore  *mongo.Client
	firebaseApp *firebase.App
}

func (a *ApiServer) New(address string,
	db *sql.DB, petDb *sql.DB, shopeStoreDb *mongo.Client, firebaseApp *firebase.App) {
	a.address = address
	a.userdb = db
	a.petsDB = petDb
	a.shopeStore = shopeStoreDb
	a.firebaseApp = firebaseApp

}

func (a *ApiServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := users.NewStore(a.userdb)
	petStore := pets.NewStore(a.petsDB)
	ShopStore := localstores.NewStore(a.shopeStore)

	pethandler := pets.NewHandler(petStore)
	userHandler := users.NewHandler(userStore)
	localStoreHandler := localstores.NewHandler(ShopStore, a.firebaseApp)

	pethandler.RegisterRoutes(subRouter)
	userHandler.RegisterRoutes(subRouter)
	localStoreHandler.RegisterRoutes(subRouter)

	log.Println("listening on port", a.address)
	handler := cors.Default().Handler(router)
	return http.ListenAndServe(a.address, handler)
}
