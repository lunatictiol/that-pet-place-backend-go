package pets

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lunatictiol/that-pet-place-backend-go/types"
)

type Handler struct {
	store types.PetStore
}

func NewHandler(store types.PetStore) *Handler {
	return &Handler{
		store: store,
	}
}
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/addPet", h.handleAddPet).Methods("POST")
	router.HandleFunc("/getPetDetails", h.handleGetPet).Methods("Get")
}

func (h *Handler) handleAddPet(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) handleGetPet(w http.ResponseWriter, r *http.Request) {

}
