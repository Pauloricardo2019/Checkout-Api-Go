package address

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

	err = db.AutoMigrate(model.Address{})
	if err != nil {
		panic(err)
	}
}

type GetByIDFn func(ID string) (*model.Address, error)

func GetByID(ID string) (*model.Address, error) {
	result := &model.Address{}
	tx := db.First(result, "id = ?", ID)

	return result, tx.Error
}

type CreateFn func(user *model.Address) error

func Create(user *model.Address) error {
	tx := db.Create(user)
	return tx.Error
}

type UpdateFn func(user *model.Address) error

func Update(user *model.Address) error {
	tx := db.Save(user)
	return tx.Error
}

type DeleteFn func(ID string) error

func Delete(ID string) error {
	tx := db.Delete(model.Address{}, "id = ?", ID)
	return tx.Error
}
