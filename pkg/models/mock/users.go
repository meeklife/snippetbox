package mock

import (
	"time"

	"meeklife.net/snippetbox/pkg/models"
)

var mockUser = &models.User{
	ID:     1,
	Name:   "Mims",
	Email:  "mims@meeklife.net",
	Create: time.Now(),
	Active: true,
}

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "mims@meeklife.net":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "mims@meeklife.net":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) ChangePassword(int, string, string) error {
	return nil
}
