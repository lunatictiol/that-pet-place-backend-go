package pets

import (
	"database/sql"
	"fmt"

	"github.com/lunatictiol/that-pet-place-backend-go/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(pet types.Pet) error {
	_, err := s.db.Exec("INSERT INTO pets (name, gender, user_id, dob, neutered, vaccinated) VALUES ($1, $2, $3, $4, $5, $6)", pet.Name, pet.Gender, pet.User_ID, pet.Dob, pet.Neutered, pet.Vaccinated)
	if err != nil {
		fmt.Print("here", err)
		return err
	}

	return nil
}

func scanUsersFromRows(row *sql.Rows) (*types.Pet, error) {
	pet := new(types.Pet)
	err := row.Scan(
		&pet.ID,
		&pet.Name,
		&pet.Gender,
		&pet.Dob,
		&pet.User_ID,
		&pet.Vaccinated,
		&pet.Neutered,
		&pet.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return pet, nil
}
