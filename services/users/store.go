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
	if us.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return us, nil

}

func (s *Store) FindUserById(id int) (*types.User, error) {
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

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES ($1, $2, $3, $4)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		fmt.Print("here", err)
		return err
	}

	return nil
}

func scanUsersFromRows(row *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}
