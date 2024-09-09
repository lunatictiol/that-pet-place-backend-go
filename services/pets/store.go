package pets

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *Store {
	return &Store{db: db}
}

// func (s *Store) CreatePet(pet types.Pet) (int64, error) {

// }

// func (s *Store) FindPetById(id int) (*types.Pet, error) {

// }
