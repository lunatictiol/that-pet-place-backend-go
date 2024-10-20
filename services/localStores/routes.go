package localstores

import (
	"context"
	"fmt"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/lunatictiol/that-pet-place-backend-go/types"
	"github.com/lunatictiol/that-pet-place-backend-go/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	store       types.ShopStore
	firebaseApp *firebase.App
}

func NewHandler(store types.ShopStore, firebaseApp *firebase.App) *Handler {
	return &Handler{
		store:       store,
		firebaseApp: firebaseApp,
	}
}
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user/getPetShops", h.GetAllStores).Methods("Get")
	router.HandleFunc("/user/getShopDetails", h.GetShopDetails).Methods("Get")
	router.HandleFunc("/user/petShopsNearUser", h.GetStoresNearUser).Methods("Get")
	router.HandleFunc("/user/bookAppointment", h.BookAnAppointment).Methods("POST")

}

func (h *Handler) GetAllStores(w http.ResponseWriter, r *http.Request) {

	ps, err := h.store.GetAllShops()
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, fmt.Errorf("error: %s", err))
		return
	}
	utils.WriteJson(w, http.StatusOK, ps)

}
func (h *Handler) GetShopDetails(w http.ResponseWriter, r *http.Request) {

	idString := r.URL.Query().Get("storeID")
	objectID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("error:%s is not a valid id", err))
		return
	}
	ps, err := h.store.GetShopDetails(objectID)
	if err != nil {
		utils.WriteJsonError(w, http.StatusInternalServerError, fmt.Errorf("error: %s", err))
		return
	}
	utils.WriteJson(w, http.StatusOK, ps)

}

func (h *Handler) GetStoresNearUser(w http.ResponseWriter, r *http.Request) {
	var payload types.PetStoreLocationUploadPayload
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
	sh, err := h.store.GetServicesNearLocation(payload.Latitude, payload.Longitute)
	if err != nil {
		utils.WriteJsonError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", err))
		fmt.Println(err)
		return
	}
	utils.WriteJson(w, http.StatusOK, sh)

}

func (h *Handler) BookAnAppointment(w http.ResponseWriter, r *http.Request) {
	var payload types.AppointmentPayload
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
	ap, err := h.store.BookAppointment(payload)
	if err != nil {
		println("hereeee")
		utils.WriteJsonError(w, http.StatusInternalServerError, err)

		return
	}
	//send store request
	utils.WriteJson(w, http.StatusOK, ap)

}
func (h *Handler) sendFirebaseNotification(storeID string, ap types.Appointment) error {
	ctx := context.Background()
	// Get Firebase messaging client
	client, err := h.firebaseApp.Messaging(ctx)
	if err != nil {
		return err
	}

	// Compose message
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "New Appointment Booked",
			Body:  fmt.Sprintf("An appointment has been booked for %s", ap.AppointmentDate),
		},
		Topic: storeID, // Assuming each store has its own topic for admin notifications
	}

	// Send notification
	_, err = client.Send(ctx, message)
	return err
}
