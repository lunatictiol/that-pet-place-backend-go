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
	_, err := s.FindPetByUserIdandName(pet.Name, pet.User_ID)
	if err == nil {
		return "", fmt.Errorf("pet already exists with that name")
	}
	_, err = s.db.Exec("INSERT INTO pets (name, gender, user_id,neutered,vaccinated,species,breed,age) VALUES ($1, $2, $3, $4,$5,$6,$7,$8)", pet.Name, pet.Gender, pet.User_ID, pet.Neutered, pet.Vaccinated, pet.Species, pet.Breed, pet.Age)
	if err != nil {

		return "", err
	}
	p, err := s.FindPetByUserIdandName(pet.Name, pet.User_ID)
	if err != nil {
		return "", err
	}

	return p.ID, nil
}
func (s *Store) UpdatePet(pet types.UpdatePet) (string, error) {
	_, err := s.db.Exec("UPDATE pets SET name=$1 , gender = $2,neutered= $3,vaccinated= $4,species= $5,breed= $6,age= $7 WHERE id =$8", pet.Name, pet.Gender, pet.Neutered, pet.Vaccinated, pet.Species, pet.Breed, pet.Age, pet.ID)
	if err != nil {

		return "", err
	}

	return pet.ID, nil
}

func (s *Store) UploadPetProfile(id string, profileUrl string) error {

	_, err := s.db.Exec("UPDATE pets SET profile = $1 WHERE id = $2 ", profileUrl, id)
	if err != nil {
		return err
	}
	return nil

}
func (s *Store) FindPetByUserIdandName(name string, id string) (*types.Pet, error) {
	rows, err := s.db.Query("SELECT * FROM pets WHERE name= $1 AND user_id = $2", name, id)
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
		return nil, fmt.Errorf("pet not found")
	}

	return p, nil

}
func (s *Store) FindPetById(id string) (*types.Pet, error) {
	rows, err := s.db.Query("SELECT * FROM pets WHERE id= $1", id)
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
		return nil, fmt.Errorf("pet not found")
	}

	return p, nil

}
func (s *Store) GetAllPets(userId string) ([]types.Pet, error) {
	rows, err := s.db.Query("SELECT * FROM pets WHERE user_id = $1", userId)
	if err != nil {
		println("WSX")
		return nil, err
	}
	defer rows.Close()

	var results []types.Pet
	for rows.Next() {
		pet := types.Pet{}
		err = rows.Scan(
			&pet.Name,
			&pet.Gender,
			&pet.Neutered,
			&pet.Vaccinated,
			&pet.Species,
			&pet.Breed,
			&pet.Profile,
			&pet.User_ID,
			&pet.ID,
			&pet.Age)
		if err != nil {
			println(err)
			return nil, err
		}

		results = append(results, pet)
	}

	if err := rows.Err(); err != nil {
		println(err)
		return nil, err
	}

	return results, nil
}
func (s *Store) DeletePet(petID string) (string, error) {
	// Execute the SQL DELETE statement
	_, err := s.db.Exec("DELETE FROM pets WHERE id = $1", petID)
	if err != nil {
		return "", err
	}

	return petID, nil
}
func scanPetsFromRows(row *sql.Rows) (*types.Pet, error) {
	pet := new(types.Pet)
	err := row.Scan(
		&pet.Name,
		&pet.Gender,
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
