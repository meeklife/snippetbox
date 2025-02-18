package mysql

import (
	"reflect"
	"testing"
	"time"

	"meeklife.net/snippetbox/pkg/models"
)

func TestUSerModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID:     1,
				Name:   "Mims",
				Email:  "mims@meeklife.net",
				Create: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active: true,
			},
			wantError: nil,
		},
		{
			name:      "Zero ID",
			userID:    0,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
		{
			name:      "non-existent ID",
			userID:    2,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{db}

			user, err := m.Get(tt.userID)

			if err != tt.wantError {
				t.Errorf("want this error: %v; got this error: %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want this user: %v; got this user: %v", tt.wantUser, user)
			}
		})
	}
}
