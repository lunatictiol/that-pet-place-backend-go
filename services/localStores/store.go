package localstores

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/lunatictiol/that-pet-place-backend-go/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *Store {

	return &Store{db: db}
}

// get all shops
func (s *Store) GetAllShops() ([]types.PetShopDetails, error) {
	collection := s.db.Database("PetServicesData").Collection("PetServices")
	result, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No appointment found with the given id
		}
		return nil, err
	}
	var ps []types.PetShopDetails
	err = result.All(context.Background(), &ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

// filterd
func (s *Store) GetAllShopsBasedOnService(filter string) ([]types.PetShopDetails, error) {
	collection := s.db.Database("PetServicesData").Collection("PetServices")

	// Find all pet shops
	result, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No pet shops found
		}
		return nil, err
	}

	var ps []types.PetShopDetails
	err = result.All(context.Background(), &ps)
	if err != nil {
		return nil, err
	}

	// Filter pet shops that offer the specific service
	var filteredShops []types.PetShopDetails
	for _, shop := range ps {
		for _, service := range shop.Services {
			if service.Name == filter {
				filteredShops = append(filteredShops, shop)
				break
			}
		}
	}

	// Return the filtered pet shops
	return filteredShops, nil

}

// shop details
func (s *Store) GetShopDetails(id primitive.ObjectID) ([]types.PetShopDetails, error) {
	collection := s.db.Database("PetServicesData").Collection("PetServices")
	result, err := collection.Find(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No appointment found with the given id
		}
		return nil, err
	}
	var ps []types.PetShopDetails
	err = result.All(context.Background(), &ps)
	if err != nil {
		return nil, err
	}
	return ps, nil

}

// near user
func (s *Store) GetServicesNearLocation(latitude float64, longitude float64) ([]types.PetShop, error) {
	collection := s.db.Database("PetServicesData").Collection("PetServices")

	// no closer than 0 meters and no farther than 10000 meters
	filter := bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{latitude, longitude},
				},
				"$maxDistance": 10000, // maximum distance in meters

			},
		},
	}

	// Find documents matching the filter using a cursor
	result, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("error finding documents: %w", err)
	}
	var ps []types.PetShop
	err = result.All(context.Background(), &ps)
	if err != nil {
		return nil, err
	}
	println(len(ps))
	for i := 0; i < len(ps); i++ {
		dis := Distance(latitude, longitude, ps[i].Location.Coordinates[0], ps[i].Location.Coordinates[1], "K")
		ps[i].Distance = math.Round(dis*100) / 100
	}
	return ps, nil

}

// book
func (s *Store) BookAppointment(ap types.AppointmentPayload) (types.Appointment, error) {
	sc := s.db.Database("PetServicesData").Collection("PetServices")
	apc := s.db.Database("PetServicesData").Collection("appointments")
	// Step 1: Find clinic by ClinicID in the sc (PetServices) collection
	var clinic struct {
		ClinicName    string `bson:"name"`
		ClinicAddress string `bson:"address"`
	}
	println(ap.CLinicID)
	clinicID, err := primitive.ObjectIDFromHex(ap.CLinicID)
	if err != nil {
		return types.Appointment{}, errors.New("invalid id")
	}

	filter := bson.M{"_id": clinicID}
	err = sc.FindOne(context.Background(), filter).Decode(&clinic)
	if err != nil {
		return types.Appointment{}, errors.New("clinic not found")
	}

	// Step 2: Create a new Appointment document
	newAppointment := types.Appointment{
		ID:                  primitive.NewObjectID(),
		DoctorName:          ap.DoctorName,
		DoctorQualification: ap.DoctorQualification,
		ClinicName:          clinic.ClinicName,
		ClinicAddress:       clinic.ClinicAddress,
		AppointmentDate:     ap.AppointmentDate,
		Status:              ap.Status,
		Price:               ap.Price,
		PetName:             ap.PetName,
		UserID:              ap.UserID,
		CLinicID:            clinicID,
		Confirmation:        "pending",
	}

	// Step 3: Insert the new appointment into the apc (appointments) collection
	_, err = apc.InsertOne(context.Background(), newAppointment)
	if err != nil {
		return types.Appointment{}, err
	}

	// Step 4: Return the new appointment
	return newAppointment, nil
}

// register clinic
func (s *Store) RegisterShop(rp types.RegisterShopPayload) (interface{}, error) {

	psa := s.db.Database("PetServicesData").Collection("PetServicesAuth")
	sp, err := psa.InsertOne(context.Background(), rp)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return sp.InsertedID, nil

}

// get user apponitments
func (s *Store) GetAllAppointments(id string) ([]types.Appointment, error) {
	collection := s.db.Database("PetServicesData").Collection("appointments")
	result, err := collection.Find(context.Background(), bson.D{{Key: "user_id", Value: id}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No appointment found with the given id
		}
		return nil, err
	}
	var ap []types.Appointment
	err = result.All(context.Background(), &ap)
	if err != nil {
		return nil, err
	}
	return ap, nil
}
func (s *Store) GetAllAppointmentsForStore(id string) ([]types.Appointment, error) {
	collection := s.db.Database("PetServicesData").Collection("appointments")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result, err := collection.Find(context.Background(), bson.D{{Key: "clinic_id", Value: objectId}})
	if err != nil {
		return nil, err
	}
	var ap []types.Appointment
	err = result.All(context.Background(), &ap)
	if err != nil {
		return nil, err
	}
	return ap, nil
}
func (s *Store) UpdateAppointmentConfirmation(id string, price int64, confirmed bool) (types.Appointment, error) {
	// Convert the string id to a MongoDB ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return types.Appointment{}, err
	}

	// Get the appointments collection
	collection := s.db.Database("PetServicesData").Collection("appointments")

	// Define the filter to find the appointment by id
	filter := bson.M{"_id": objectId}

	// Find the appointment
	var appointment types.Appointment
	err = collection.FindOne(context.Background(), filter).Decode(&appointment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.Appointment{}, nil // No appointment found with the given id
		}
		return types.Appointment{}, err
	}

	// Update the "confirmed" field of the appointment
	update := bson.M{
		"$set": bson.M{"confirmed": confirmed, "confirmation": "done", "price": price},
	}

	// Update the document
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return types.Appointment{}, err
	}

	// Optionally, you can return the updated appointment by fetching it again
	err = collection.FindOne(context.Background(), filter).Decode(&appointment)
	if err != nil {
		return types.Appointment{}, err
	}

	return appointment, nil
}

func (s *Store) UpdateAppointmentStatus(id string, status string) (types.Appointment, error) {
	// Convert the string id to a MongoDB ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return types.Appointment{}, err
	}

	// Get the appointments collection
	collection := s.db.Database("PetServicesData").Collection("appointments")

	// Define the filter to find the appointment by id
	filter := bson.M{"_id": objectId}

	// Find the appointment
	var appointment types.Appointment
	err = collection.FindOne(context.Background(), filter).Decode(&appointment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.Appointment{}, nil // No appointment found with the given id
		}
		return types.Appointment{}, err
	}

	// Update the "confirmed" field of the appointment
	update := bson.M{
		"$set": bson.M{"status": status},
	}

	// Update the document
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return types.Appointment{}, err
	}

	// Optionally, you can return the updated appointment by fetching it again
	err = collection.FindOne(context.Background(), filter).Decode(&appointment)
	if err != nil {
		return types.Appointment{}, err
	}

	return appointment, nil
}

//send ratings

func (s *Store) UpdateShopRating(petShopId string, newRating float64) error {
	// Convert string ID to ObjectID
	objectId, err := primitive.ObjectIDFromHex(petShopId)
	if err != nil {
		return err
	}

	// Get the pet shops collection
	collection := s.db.Database("PetServicesData").Collection("PetServices")

	// Find the current pet shop details by ID
	var petShop types.PetShopDetails
	filter := bson.M{"_id": objectId}
	err = collection.FindOne(context.Background(), filter).Decode(&petShop)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("no pet shop found with the given ID")
		}
		return err
	}

	// Calculate new average rating
	oldRating := petShop.Ratings
	oldRatingCount := petShop.RatingCount
	newAverageRating := ((oldRating * float64(oldRatingCount)) + newRating) / float64(oldRatingCount+1)

	// Increment the rating count
	newRatingCount := oldRatingCount + 1

	// Update the pet shop with new rating and rating count
	update := bson.M{
		"$set": bson.M{
			"ratings":      newAverageRating,
			"rating_count": newRatingCount,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	// Optionally, retrieve the updated pet shop details
	err = collection.FindOne(context.Background(), filter).Decode(&petShop)
	if err != nil {
		return err
	}

	return nil
}

//login clinic

//add clinic details

//add doctor

//add product

//add clicnic profile photo
