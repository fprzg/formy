package models

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestUsersInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := SetupTestDB(t)

	tests := []struct {
		TestName      string
		name          string
		email         string
		password      string
		expectedError error
	}{
		{
			TestName:      "Successful insertion",
			name:          "Bob",
			email:         "bob@example.com",
			password:      "pass123",
			expectedError: nil,
		},
		{
			TestName:      "Duplicated email",
			name:          "Alice",
			email:         ValidUserEmail,
			password:      "pass123",
			expectedError: ErrDuplicateEmail,
		},
		{
			TestName:      "Empty name",
			name:          "",
			email:         "emptyname@example.com",
			password:      "pass123",
			expectedError: ErrInvalidInput,
		},
		{
			TestName:      "Invalid email",
			name:          "Invalid Email",
			email:         "not-an-email",
			password:      "pass123",
			expectedError: ErrInvalidInput,
		},
		{
			TestName:      "Empty password",
			name:          "No Password",
			email:         "nopass@example.com",
			password:      "",
			expectedError: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			err := m.Users.Insert(tt.name, tt.email, tt.password)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedError, err.Error())
			}
		})
	}
}

func TestUsersAuthenticate(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := SetupTestDB(t)

	tests := []struct {
		TestName      string
		email         string
		password      string
		wantID        int
		expectedError error
	}{
		{
			TestName:      "Successful authentication",
			email:         ValidUserEmail,
			password:      ValidUserPassword,
			wantID:        1,
			expectedError: nil,
		},
		{
			TestName:      "Wrong email",
			email:         "bob@example.com",
			password:      ValidUserPassword,
			wantID:        0,
			expectedError: ErrInvalidCredentials,
		},
		{
			TestName:      "Wrong password",
			email:         ValidUserEmail,
			password:      "wrongpass",
			wantID:        0,
			expectedError: ErrInvalidCredentials,
		},
		{
			TestName:      "Empty email",
			email:         "",
			password:      "pass123",
			wantID:        0,
			expectedError: ErrInvalidCredentials,
		},
		{
			TestName:      "Empty password",
			email:         ValidUserEmail,
			password:      "",
			wantID:        0,
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			userID, err := m.Users.Authenticate(tt.email, tt.password)
			if tt.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, userID)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedError, err.Error())
			}
		})
	}
}

func TestUsersExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := SetupTestDB(t)

	tests := []struct {
		TestName       string
		userID         int
		expectedExists bool
	}{
		{
			TestName:       "Valid ID",
			userID:         1,
			expectedExists: true,
		},
		{
			TestName:       "Zero ID",
			userID:         0,
			expectedExists: false,
		},
		{
			TestName:       "Negative ID",
			userID:         -1,
			expectedExists: false,
		},
		{
			TestName:       "Non-existent ID",
			userID:         999,
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			exists, err := m.Users.Exists(tt.userID)
			assert.Equal(t, tt.expectedExists, exists)
			assert.NoError(t, err)
		})
	}
}

func TestUsersGet(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := SetupTestDB(t)

	tests := []struct {
		TestName      string
		id            int
		name          string
		email         string
		expectedError error
	}{
		{
			TestName:      "Valid user ID",
			id:            1,
			name:          ValidUserName,
			email:         ValidUserEmail,
			expectedError: nil,
		},
		{
			TestName:      "Zero ID",
			id:            0,
			expectedError: ErrNoRecord,
		},
		{
			TestName:      "Negative ID",
			id:            -1,
			expectedError: ErrNoRecord,
		},
		{
			TestName:      "Non-existent ID",
			id:            999,
			expectedError: ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			user, err := m.Users.Get(tt.id)
			if tt.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.name, user.Name)
				assert.Equal(t, tt.email, user.Email)
			} else {
				assert.EqualError(t, tt.expectedError, err.Error())
			}
		})
	}
}

func TestUsersUpdateName(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := SetupTestDB(t)

	id := InsertTestUser(t, m, "Dave", "dave@example.com", "pass")

	tests := []struct {
		TestName     string
		newName      string
		currentPass  string
		expectError  error
		expectedName string
	}{
		{
			TestName:     "Successful update",
			newName:      "Dave Updated",
			currentPass:  "pass",
			expectError:  nil,
			expectedName: "Dave Updated",
		},
		{
			TestName:     "Wrong password",
			newName:      "Dave WrongPass",
			currentPass:  "wrong",
			expectError:  ErrInvalidCredentials,
			expectedName: "Dave Updated",
		},
		{
			TestName:     "Empty name",
			newName:      "",
			currentPass:  "pass",
			expectError:  ErrInvalidInput,
			expectedName: "Dave Updated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			err := m.Users.UpdateName(id, tt.newName, tt.currentPass)
			if tt.expectError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.expectError)
			}

			user, err := m.Users.Get(id)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedName, user.Name)
		})
	}
}

func TestUsersUpdateEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := SetupTestDB(t)

	id := InsertTestUser(t, m, "Eve", "eve@old.com", "pass")
	_ = InsertTestUser(t, m, "Other", "used@example.com", "pass")

	tests := []struct {
		TestName     string
		userID       int
		newEmail     string
		currentPass  string
		expectError  error
		expectedAuth string
	}{
		{
			TestName:     "Successful update",
			userID:       id,
			newEmail:     "eve@new.com",
			currentPass:  "pass",
			expectError:  nil,
			expectedAuth: "eve@new.com",
		},
		{
			TestName:     "Wrong password",
			userID:       id,
			newEmail:     "eve@wrongpass.com",
			currentPass:  "wrong",
			expectError:  ErrInvalidCredentials,
			expectedAuth: "eve@new.com",
		},
		{
			TestName:     "Duplicated email",
			userID:       id,
			newEmail:     "used@example.com",
			currentPass:  "pass",
			expectError:  ErrDuplicateEmail,
			expectedAuth: "eve@new.com",
		},
		{
			TestName:     "Invalid ID",
			userID:       0,
			newEmail:     "used@example.com",
			currentPass:  "pass",
			expectError:  ErrInvalidCredentials,
			expectedAuth: "eve@new.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			err := m.Users.UpdateEmail(tt.userID, tt.newEmail, tt.currentPass)
			if tt.expectError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.expectError)
			}

			_, err = m.Users.Authenticate(tt.expectedAuth, "pass")
			assert.NoError(t, err)
		})
	}
}

func TestUsersUpdatePassword(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := SetupTestDB(t)

	id := InsertTestUser(t, m, "Frank", "frank@example.com", "oldpass")

	tests := []struct {
		TestName        string
		currentPass     string
		newPass         string
		expectError     error
		canLoginWith    string
		cannotLoginWith string
	}{
		{
			TestName:        "Successful update",
			currentPass:     "oldpass",
			newPass:         "newpass",
			expectError:     nil,
			canLoginWith:    "newpass",
			cannotLoginWith: "oldpass",
		},
		{
			TestName:        "Wrong current password",
			currentPass:     "wrong",
			newPass:         "another",
			expectError:     ErrInvalidCredentials,
			canLoginWith:    "newpass",
			cannotLoginWith: "another",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			err := m.Users.UpdatePassword(id, tt.currentPass, tt.newPass)
			if tt.expectError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.expectError)
			}

			_, err = m.Users.Authenticate("frank@example.com", tt.canLoginWith)
			assert.NoError(t, err)

			if tt.cannotLoginWith != "" {
				_, err = m.Users.Authenticate("frank@example.com", tt.cannotLoginWith)
				assert.ErrorIs(t, err, ErrInvalidCredentials)
			}
		})
	}
}
