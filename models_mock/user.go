package models_mock

import (
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/xssnick/crawlyzer-auth/models"
)

type MUserStore struct {
	FakeError error
}

var TestUUID = uuid.FromStringOrNil("4f477db2-17b7-432d-ad0e-c2a098cbd2b0")

func (us *MUserStore) Create(email, password string) (uuid.UUID, error) {
	if us.FakeError != nil {
		return uuid.Nil, us.FakeError
	}
	return TestUUID, nil
}

func (us *MUserStore) Auth(sesid string) (uuid.UUID, error) {
	if us.FakeError != nil {
		return uuid.Nil, us.FakeError
	}

	return TestUUID, nil
}

func (us *MUserStore) Logout(sesid string) error {
	if us.FakeError != nil {
		return us.FakeError
	}

	return nil
}

func (us *MUserStore) Login(email, password string) (string, error) {
	if us.FakeError != nil {
		return "", us.FakeError
	}
	return "6e536fff-baaf-4ca7-a067-352bafeb6ee3", nil
}

func (us *MUserStore) GetAll() ([]models.User, error) {
	if us.FakeError != nil {
		return nil, us.FakeError
	}

	return []models.User{
		{
			ID:    uuid.FromStringOrNil("4f477db2-17b7-432d-ad0e-c2a098cbd2b0"),
			Email: "tester@exter.com",
		}, {
			ID:    uuid.FromStringOrNil("6e536fff-bcaf-4ca9-a067-352bafeb6ed2"),
			Email: "western@eastern.so",
		}, {
			ID:    uuid.FromStringOrNil("d96bee74-07c5-40ca-b0cc-c0e04d4a7589"),
			Email: "chester@pepster.net",
		},
	}, nil
}
