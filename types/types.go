package types

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterUserPayload struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=5,max=8" `
}
type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required" `
}
type PetPayload struct {
	Name       string `json:"name" validate:"required"`
	Gender     string `json:"gender" validate:"required"`
	User_ID    string `json:"user_id" validate:"required"`
	Dob        string `json:"dob" validate:"required"`
	Neutered   bool   `json:"neutered" validate:"required"`
	Vaccinated bool   `json:"vaccinated" validate:"required"`
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
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
	User_ID    string    `json:"user_id"`
	Dob        string    `json:"dob"`
	Neutered   bool      `json:"neutered"`
	Breed      string    `json:"breed"`
	Species    string    `json:"species"`
	Vaccinated bool      `json:"vaccinated"`
	Age        int       `json:"age"`
	Profile    string    `json:"profile"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserStore interface {
	FindUserByEmail(email string) (*User, error)
	FindUserById(id string) (*User, error)
	CreateUser(User) (string, error)
	UploadProfile(id string, profileUrl string) error
}

type PetStore interface {
	FindPetById(id string) (*Pet, error)
	CreatePet(Pet) (string, error)
}

type Manager struct {
	Connection *mongo.Client
	Ctx        context.Context
	Cancel     context.CancelFunc
}
