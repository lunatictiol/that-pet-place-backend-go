package pets

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lunatictiol/that-pet-place-backend-go/types"
	"github.com/lunatictiol/that-pet-place-backend-go/utils"
	"google.golang.org/api/option"
)

type Handler struct {
	store types.PetStore
}

const (
	projectID  = "thatpetplace"
	bucketName = "pet-profile-bucket"
	//local
	//credentials = "./application_default_credentials.json"
	//remote
	credentials = "/etc/secrets/application_default_credentials.json"
)

func NewHandler(store types.PetStore) *Handler {
	return &Handler{
		store: store,
	}
}
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/addPet", h.handleAddPet).Methods("POST")
	router.HandleFunc("/getPetDetails", h.handleGetAllPets).Methods("Get")
	router.HandleFunc("/getAllPets", h.handleGetAllPets).Methods("Get")
	router.HandleFunc("/uploadPetProfile", h.handleProfileUpload).Methods("POST")

}

func (h *Handler) handleGetAllPets(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userID")
	uid, err := uuid.Parse(userId)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, fmt.Errorf("error parsing id: %s", userId))

		return
	}

	p, err := h.store.GetAllPets(uid)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, fmt.Errorf("error retrieveing data of id: %s", userId))

		return
	}
	if p == nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("no pets found for id: %s", userId))
		return

	}
	utils.WriteJson(w, http.StatusCreated, map[string]any{"message": "pets retrieved successfully", "pets": p})

}
func (h *Handler) handleAddPet(w http.ResponseWriter, r *http.Request) {
	var payload types.PetPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, err)
		fmt.Println(err)
		return
	}
	if err := utils.Validator.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", error))
		fmt.Println(err)
		return
	}
	uid, err := uuid.Parse(payload.User_ID)
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid user id %v", err))
		fmt.Println(err)
		return
	}

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
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]any{"message": "pet added successful", "id": uId})

}
func (h *Handler) handleProfileUpload(w http.ResponseWriter, r *http.Request) {
	uId := r.URL.Query().Get("petID")
	_, err := h.store.FindPetById(uId)

	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("user doesn't exists with id %s", uId))
		println(err)
		return
	}

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

	err = h.store.UploadPetProfile(uId, u.Path)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		println(err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "upload successful", "url": u.Path})

}
