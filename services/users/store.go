package users

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

func (s *Store) FindUserByEmail(email string) (*types.User, error) {

	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	us := new(types.User)
	for rows.Next() {
		us, err = scanUsersFromRows(rows)
		if err != nil {
			return nil, err
		}

	}
	if us.ID == "" {
		return nil, fmt.Errorf("user not found")
	}
	return us, nil
}

func (s *Store) FindUserById(id string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanUsersFromRows(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == "" {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) CreateUser(user types.User) (string, error) {
	_, err := s.db.Exec("INSERT INTO users (first_name, last_name, email, password,phone_number) VALUES ($1, $2, $3, $4,$5)", user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNumber)
	if err != nil {

		return "", err
	}
	u, err := s.FindUserByEmail(user.Email)
	if err != nil {
		return "", err
	}
	println(u.ID)
	return u.ID, nil
}
func (s *Store) UploadProfile(id string, profileUrl string) error {

	_, err := s.db.Exec("UPDATE users SET profile = $1 WHERE id = $2 ", profileUrl, id)
	if err != nil {
		return err
	}
	return nil

}

func scanUsersFromRows(row *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := row.Scan(
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.Profile,
		&user.PhoneNumber,
		&user.ID,
		&user.PetID,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}
