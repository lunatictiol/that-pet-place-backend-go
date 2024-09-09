package users

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/lunatictiol/that-pet-place-backend-go/services/auth"
	"github.com/lunatictiol/that-pet-place-backend-go/types"
	"github.com/lunatictiol/that-pet-place-backend-go/utils"
	"google.golang.org/api/option"
)

type Handler struct {
	store types.UserStore
}

const (
	projectID  = "thatpetplace"
	bucketName = "pet-parents-profile"
	//local
	//credentials = "./application_default_credentials.json"
	credentials = "./etc/secrets/application_default_credentials.json"
)

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/uploadProfile", h.handleProfileUpload).Methods("POST")
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

	utils.WriteJson(w, http.StatusCreated, map[string]any{"message": "Login successful", "token": token, "userId": u.ID})

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

	uId, err := h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]any{"message": "Registeration successful", "id": uId})

}
func (h *Handler) handleProfileUpload(w http.ResponseWriter, r *http.Request) {

	// Parse the uploaded file
	f, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		println(err)
		return
	}
	defer f.Close()
	// Create a context with a timeout
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Initialize the storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentials))
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		println(err)
		return
	}
	// Create a unique object name for the uploaded file
	objectName := fmt.Sprintf("%s/%s", "photos", handler.Filename)

	// Upload the file to the bucket
	sw := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.Copy(sw, f); err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		println(err)
		return
	}
	if err := sw.Close(); err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		println(err)
		return
	}

	u, err := url.Parse("/" + bucketName + "/" + sw.Attrs().Name)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		println(err)
		return
	}

	uId := r.URL.Query().Get("userID")
	userId, err := strconv.Atoi(uId)
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("inalid id %s", string(uId)))
		println(err)
		return

	}

	_, err = h.store.FindUserById(int(userId))

	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("user doesn't exists with id %s", uId))
		println(err)
		return
	}
	err = h.store.UploadProfile(userId, u.Path)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		println(err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "upload successful", "url": u.Path})

}
