package pets

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lunatictiol/that-pet-place-backend-go/types"
	"github.com/lunatictiol/that-pet-place-backend-go/utils"
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
	var payload types.PetPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validator.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", error))
		return
	}
	uid, err := uuid.Parse(payload.User_ID)
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid user id %v", err))
		return
	}
	println(uid.String())
	uId, err := h.store.CreatePet(types.Pet{
		Name:       payload.Name,
		Gender:     payload.Gender,
		User_ID:    uid,
		Dob:        payload.Dob,
		Neutered:   payload.Neutered,
		Breed:      payload.Breed,
		Species:    payload.Species,
		Age:        payload.Age,
		Vaccinated: payload.Vaccinated,
	})
	if err != nil {
		fmt.Printf("here")
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		fmt.Printf(err.Error())
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]any{"message": "pet added successful", "id": uId})

}
func (h *Handler) handleGetPet(w http.ResponseWriter, r *http.Request) {

}
