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

func (s *Store) GetAllShops() ([]types.PetShop, error) {
	collection := s.db.Database("PetServicesData").Collection("PetServices")
	result, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var ps []types.PetShop
	err = result.All(context.Background(), &ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}
func (s *Store) GetShopDetails(id primitive.ObjectID) ([]types.PetShopDetails, error) {
	collection := s.db.Database("PetServicesData").Collection("PetServices")
	result, err := collection.Find(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return nil, err
	}
	var ps []types.PetShopDetails
	err = result.All(context.Background(), &ps)
	if err != nil {
		return nil, err
	}
	return ps, nil

}
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

//register clinic
// func (s *Store) RegisterShop(ap types.RegisterShopPayload) (types.Appointment, error) {
// // 	sc := s.db.Database("PetServicesData").Collection("PetServices")
// // 	apc := s.db.Database("PetServicesData").Collection("appointments")
// // }
//login clinic

//add clinic details
//add doctor
//add product
//approve appointment
//send upcomming appointments
//send ratings
//send past appointments
//cancel appointment
//complete appointment
