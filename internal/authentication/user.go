package authentication

import "github.com/google/uuid"

type User struct {
	ID string
}

func (user *User) SetID(userID string) {
	user.ID = userID
}

func (user *User) generateUniqueUserID() {
	UUID := uuid.New()
	user.SetID(UUID.String())
}

func (user *User) GetID() string {
	if user.ID == "" {
		user.generateUniqueUserID()
	}

	return user.ID
}
