package types

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterUserPayload struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=5,max=16" `
}
type UpdateUserPayload struct {
	Id          string `json:"id" validate:"required" `
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required" `
}

type PetPayload struct {
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	User_ID    string `json:"user_id"`
	Dob        string `json:"dob"`
	Neutered   bool   `json:"neutered"`
	Vaccinated bool   `json:"vaccinated"`
	Breed      string `json:"breed"`
	Species    string `json:"species"`
	Age        int    `json:"age"`
}
type ProfileUploadPayload struct {
	Profile string `json:"profile" validate:"required"`
	ID      string `json:"id" validate:"required"`
}
type PetProfileUploadPayload struct {
	Profile string `json:"profile" validate:"required"`
	ID      string `json:"id" validate:"required"`
}
type PetStoreLocationUploadPayload struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitute float64 `json:"longitude" validate:"required"`
}

type User struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	PetID       any       `json:"petid"`
	PhoneNumber string    `json:"phone_number"`
	Profile     any       `json:"profile"`
	CreatedAt   time.Time `json:"created_at"`
}
type Pet struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	User_ID    string `json:"user_id"`
	Dob        string `json:"dob"`
	Neutered   bool   `json:"neutered"`
	Breed      string `json:"breed"`
	Species    string `json:"species"`
	Vaccinated bool   `json:"vaccinated"`
	Age        int    `json:"age"`
	Profile    any    `json:"profile"`
}
type UpdatePet struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	Dob        string `json:"dob"`
	Neutered   bool   `json:"neutered"`
	Breed      string `json:"breed"`
	Species    string `json:"species"`
	Vaccinated bool   `json:"vaccinated"`
	Age        int    `json:"age"`
}

type RegisterShopPayload struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=5,max=16" `
}

type PetShop struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Tagline     string             `json:"tagline"`
	Ratings     float64            `json:"ratings"`
	Type        string             `json:"type"`
	Services    struct {
		Grooming        int    `json:"grooming"`
		VetinaryCheckup int    `json:"vetinary_checkup"`
		PetBoarding     int    `json:"pet_boarding"`
		Adoption        int    `json:"adoption"`
		DogWalking      int    `json:"dog_walking"`
		Training        int    `json:"training"`
		PetTaxi         int    `json:"pet_taxi"`
		Vaccination     int    `json:"vaccination"`
		Other           string `json:"other"`
	} `json:"services"`
	Location struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"location"`
	Distance float64
}

type PetShopDetails struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Tagline     string             `json:"tagline"`
	Ratings     float64            `json:"ratings"`
	Type        string             `json:"type"`
	Services    struct {
		Grooming        int    `json:"grooming"`
		VetinaryCheckup int    `json:"vetinary_checkup"`
		PetBoarding     int    `json:"pet_boarding"`
		Adoption        int    `json:"adoption"`
		DogWalking      int    `json:"dog_walking"`
		Training        int    `json:"training"`
		PetTaxi         int    `json:"pet_taxi"`
		Vaccination     int    `json:"vaccination"`
		Other           string `json:"other"`
	} `json:"services"`
	Location struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"location"`
	Distance    float64
	Doctors     []Doctor  `bson:"doctors"`
	Address     string    `bson:"address"`
	PhoneNumber string    `bson:"phone_number"`
	Profile     string    `bson:"profile"`
	Products    []Product `bson:"products"`
}
type Appointment struct {
	ID                  primitive.ObjectID `bson:"_id"`
	DoctorName          string             `bson:"doctor_name"`
	DoctorQualification string             `bson:"doctor_qualification"`
	ClinicName          string             `bson:"clinic_name"`
	ClinicAddress       string             `bson:"clinic_address"`
	AppointmentDate     time.Time          `bson:"appointment_date"`
	Status              string             `bson:"status"`
	Price               float64            `bson:"price"`
	PetName             string             `bson:"pet_name"`
	UserID              string             `bson:"user_id"`
	CLinicID            primitive.ObjectID `bson:"clinic_id"`
	Confirmation        string             `bson:"confirmation"`
	Confirmed           bool               `bson:"confirmed"`
}
type AppointmentClicnicApproval struct {
	ID    primitive.ObjectID `bson:"_id"`
	Price float64            `bson:"price"`
}
type AppointmentPayload struct {
	DoctorName          string    `json:"doctor_name"`
	AppointmentDate     time.Time `json:"appointment_date"`
	PetName             string    `json:"pet_name"`
	UserID              string    `json:"user_id"`
	CLinicID            string    `json:"clinic_id"`
	Status              string    `json:"status"`
	DoctorQualification string    `bson:"doctor_qualification"`
	Price               float64   `bson:"price"`
}

type Doctor struct {
	Name              string   `bson:"name"`
	Qualification     string   `bson:"qualification"`
	Fees              float64  `bson:"fees"`                // Fees is stored as a float
	AvailableDays     []string `bson:"available_days"`      // Available days is an array of strings
	YearsOfExperience float64  `bson:"years_of_experience"` // Years of experience is stored as a float
	Profile           string   `bson:"profile"`
}

type Product struct {
	Name        string  `bson:"name"`
	Photo       string  `bson:"photo"`
	Price       float64 `bson:"price"` // Price is stored as a float
	Description string  `bson:"description"`
}
type UserStore interface {
	FindUserByEmail(email string) (*User, error)
	FindUserById(id string) (*User, error)
	CreateUser(User) (string, error)
	UpdateUser(UpdateUserPayload) (string, error)
	UploadProfile(id string, profileUrl string) error
}

type PetStore interface {
	FindPetByUserIdandName(name string, id string) (*Pet, error)
	CreatePet(Pet) (string, error)
	UpdatePet(UpdatePet) (string, error)
	GetAllPets(userId string) ([]Pet, error)
	FindPetById(id string) (*Pet, error)
	UploadPetProfile(id string, profileUrl string) error
}
type ShopStore interface {
	GetAllShops() ([]PetShop, error)
	GetServicesNearLocation(latitude float64, longitude float64) ([]PetShop, error)
	GetShopDetails(id primitive.ObjectID) ([]PetShopDetails, error)
	BookAppointment(ap AppointmentPayload) (Appointment, error)
}

type Manager struct {
	Connection *mongo.Client
	Ctx        context.Context
	Cancel     context.CancelFunc
}
