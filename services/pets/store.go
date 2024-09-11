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

func (s *Store) CreatePet(pet types.Pet) (string, error) {
	_, err := s.db.Exec("INSERT INTO pets (name, gender, user_id, dob,neutered,vaccinated,species,breed,age) VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9)", pet.Name, pet.Gender, pet.User_ID, pet.Dob, pet.Neutered, pet.Vaccinated, pet.Species, pet.Breed, pet.Age)
	if err != nil {

		return "", err
	}
	p, err := s.FindPetById(pet.ID)
	if err != nil {
		return "", err
	}
	println(p.ID)
	return p.ID, nil
}

func (s *Store) FindPetById(id string) (*types.Pet, error) {
	rows, err := s.db.Query("SELECT * FROM pets WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	p := new(types.Pet)
	for rows.Next() {
		p, err = scanPetsFromRows(rows)
		if err != nil {
			return nil, err
		}
	}

	if p.ID == "" {
		return nil, fmt.Errorf("user not found")
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
		&pet.CreatedAt,
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
