package customer

import (
	"ravxcheckout/src/adapter/database/sql"
	model "ravxcheckout/src/internal/model/db"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {

	database, err := sql.GetGormDB()
	if err != nil {
		return
	}
	db = database

	err = db.AutoMigrate(model.Customer{})
	if err != nil {
		panic(err)
	}
}

type GetByIDFn func(ID string) (*model.Customer, error)

func GetByID(ID string) (*model.Customer, error) {
	result := &model.Customer{}
	tx := db.First(result, "id = ?", ID)

	return result, tx.Error
}

type CreateFn func(user *model.Customer) error

func Create(user *model.Customer) error {
	tx := db.Create(user)
	return tx.Error
}

type UpdateFn func(user *model.Customer) error

func Update(user *model.Customer) error {
	tx := db.Save(user)
	return tx.Error
}

type DeleteFn func(ID string) error

func Delete(ID string) error {
	tx := db.Delete(model.Customer{}, "id = ?", ID)
	return tx.Error
}
