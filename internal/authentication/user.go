package authentication

import "github.com/google/uuid"

type UserStruct struct {
	ID string
}

func (user *UserStruct) SetID(userID string) {
	user.ID = userID
}

func (user *UserStruct) generateUniqueUserID() {
	UUID, _ := uuid.NewRandom()
	user.SetID(UUID.String())
}

func (user *UserStruct) GetID() string {
	if user.ID == "" {
		user.generateUniqueUserID()
	}

	return user.ID
}

var user = &UserStruct{}

func User() *UserStruct {
	return user
}
