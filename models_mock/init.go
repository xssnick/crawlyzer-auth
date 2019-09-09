package models_mock

import "github.com/xssnick/crawlyzer-auth/models"

func InitMockStore() *models.DataStore {
	return &models.DataStore{
		User: &MUserStore{},
	}
}
