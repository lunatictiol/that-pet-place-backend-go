package users

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/lunatictiol/that-pet-place-backend-go/services/auth"
	"github.com/lunatictiol/that-pet-place-backend-go/types"
	"github.com/lunatictiol/that-pet-place-backend-go/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload
	if err := utils.ParseJson(r, &user); err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validator.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err := h.store.FindUserByEmail(user.Email)
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ValidatePassword(user.Password, u.Password) {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid password"))
		return
	}

	token, err := auth.GenerateToken(int64(u.ID))
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "Login successful", "token": token})

}
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", error))
		return
	}
	_, err := h.store.FindUserByEmail(payload.Email)

	if err == nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("user already exists with email %s", payload.Email))
		return
	}
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, "user created")

}
