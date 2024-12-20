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
	router.HandleFunc("/pet/addPet", h.handleAddPet).Methods("POST")
	router.HandleFunc("/pet/getPetDetails", h.handleGetPetDetails).Methods("Get")
	router.HandleFunc("/pet/getAllPets", h.handleGetAllPets).Methods("Get")
	router.HandleFunc("/pet/updatePet", h.handleUpdatePet).Methods("POST")
	router.HandleFunc("/pet/Delete", h.handleDeletePet).Methods("DELETE")
	router.HandleFunc("/pet/uploadPetProfile", h.handleProfileUpload).Methods("POST")

}

func (h *Handler) handleGetPetDetails(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("petID")

	p, err := h.store.FindPetById(userId)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, fmt.Errorf("error retrieveing data of id: %s", userId))

		return
	}
	if p == nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("no pets found for id: %s", userId))
		return

	}
	utils.WriteJson(w, http.StatusCreated, map[string]any{"message": "pet retrieved successfully", "pet": p})

}

func (h *Handler) handleGetAllPets(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userID")
	println(userId)
	p, err := h.store.GetAllPets(userId)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, fmt.Errorf("error retrieveing data of user id: %s", userId))

		return
	}
	if p == nil {
		utils.WriteJson(w, http.StatusOK, map[string]any{"message": "no pets", "pets": p})
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

	uId, err := h.store.CreatePet(types.Pet{
		Name:       payload.Name,
		Gender:     payload.Gender,
		User_ID:    payload.User_ID,
		Neutered:   payload.Neutered,
		Breed:      payload.Breed,
		Species:    payload.Species,
		Age:        payload.Age,
		Vaccinated: payload.Vaccinated,
	})
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, err)
		fmt.Println(err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]any{"message": "pet added successful", "id": uId})

}
func (h *Handler) handleUpdatePet(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdatePet
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

	uId, err := h.store.UpdatePet(payload)
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, err)
		fmt.Println(err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]any{"message": "pet updated successful", "id": uId})

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
	publicURL := fmt.Sprintf("https://storage.cloud.google.com/%s/photos/%s", bucketName, handler.Filename)

	u, err := url.Parse(publicURL)
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
	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "upload successful", "url": publicURL})

}

func (h *Handler) handleDeletePet(w http.ResponseWriter, r *http.Request) {
	uId := r.URL.Query().Get("petID")
	_, err := h.store.FindPetById(uId)

	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("pet doesn't exists with id %s", uId))
		println(err)
		return
	}

	str, err := h.store.DeletePet(uId)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, err)
		println(err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "pet deleted successful", "id": str})

}
