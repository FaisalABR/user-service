package seeders

import (
	"user-service/domain/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RunRoleSeeder(db *gorm.DB) {

	roles := []models.Role{
		{
			Code: "ADMIN",
			Name: "Admin",
		},
		{
			Code: "USER",
			Name: "User",
		},
	}

	for _, role := range roles {
		err := db.FirstOrCreate(&role, models.Role{Code: role.Code}).Error
		if err != nil {
			logrus.Errorf("failed to seed roles: %v", err)
			panic(err)
		}
		logrus.Infof("%s roles successfuly seeded!", role.Code)
	}

}
