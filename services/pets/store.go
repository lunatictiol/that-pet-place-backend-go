package pets

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lunatictiol/that-pet-place-backend-go/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreatePet(pet types.Pet) (uuid.UUID, error) {
	_, err := s.FindPetByUserIdandName(pet.Name, pet.User_ID)
	if err == nil {
		return uuid.Nil, fmt.Errorf("pet already exists with that name")
	}
	_, err = s.db.Exec("INSERT INTO pets (name, gender, user_id, dob,neutered,vaccinated,species,breed,age) VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9)", pet.Name, pet.Gender, pet.User_ID, pet.Dob, pet.Neutered, pet.Vaccinated, pet.Species, pet.Breed, pet.Age)
	if err != nil {
		println("1")
		return uuid.Nil, err
	}
	p, err := s.FindPetByUserIdandName(pet.Name, pet.User_ID)
	if err != nil {
		return uuid.Nil, err
	}

	return p.ID, nil
}

func (s *Store) FindPetByUserIdandName(name string, id uuid.UUID) (*types.Pet, error) {
	rows, err := s.db.Query("SELECT * FROM pets WHERE name= $1 AND user_id = $2", name, id.String())
	if err != nil {
		println("3")
		return nil, err
	}

	p := new(types.Pet)
	for rows.Next() {
		p, err = scanPetsFromRows(rows)
		if err != nil {
			return nil, err
		}
	}

	if p.ID == uuid.Nil {
		return nil, fmt.Errorf("pet not found")
	}

	return p, nil

}

func scanPetsFromRows(row *sql.Rows) (*types.Pet, error) {
	pet := new(types.Pet)
	err := row.Scan(
		&pet.Name,
		&pet.Gender,
		&pet.Dob,
		&pet.Neutered,
		&pet.Vaccinated,
		&pet.Species,
		&pet.Breed,
		&pet.Profile,
		&pet.User_ID,
		&pet.ID,
		&pet.Age,
	)

	if err != nil {
		return nil, err
	}
	return pet, nil
}
