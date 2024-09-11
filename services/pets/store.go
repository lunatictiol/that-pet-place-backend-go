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
func (s *Store) GetAllPets(userId uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := s.db.Query("SELECT * FROM pets WHERE user_id = $1", userId.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var results []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		result := make(map[string]interface{})
		for i := range columns {
			result[columns[i]] = values[i]
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
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
