package sql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetGormDB() (*gorm.DB, error) {
	dsn := "host=localhost user=user password=password dbname=ravxcheckout port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{FullSaveAssociations: true})
}
