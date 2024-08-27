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

func (s *Store) CreatePet(pet types.Pet) (int64, error) {
	p, err := s.db.Exec("INSERT INTO pets (name, gender, user_id, dob, neutered, vaccinated) VALUES ($1, $2, $3, $4, $5, $6)", pet.Name, pet.Gender, pet.User_ID, pet.Dob, pet.Neutered, pet.Vaccinated)
	if err != nil {

		return 0, err
	}
	pId, err := p.LastInsertId()
	if err != nil {

		return 0, err
	}

	return pId, nil
}

func (s *Store) FindPetById(id int) (*types.Pet, error) {
	rows, err := s.db.Query("SELECT * FROM pets WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	p := new(types.Pet)
	for rows.Next() {
		p, err = scanUsersFromRows(rows)
		if err != nil {
			return nil, err
		}
	}

	if p.ID == 0 {
		return nil, fmt.Errorf("pet not found")
	}

	return p, nil
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
