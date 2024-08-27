package types

import "time"

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=5,max=8" `
}
type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required" `
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}
type Pet struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
	User_ID    string    `json:"user_id"`
	Dob        string    `json:"dob"`
	Neutered   bool      `json:"neutered"`
	Vaccinated bool      `json:"vaccinated"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserStore interface {
	FindUserByEmail(email string) (*User, error)
	FindUserById(id int) (*User, error)
	GetUserId(email string) (int, error)
	CreateUser(User) error
}

type PetStore interface {
	FindPetById(id int) (*Pet, error)
	CreatePet(Pet) error
}
