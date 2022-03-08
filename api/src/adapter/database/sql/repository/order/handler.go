package order

import (
	"ravxcheckout/src/adapter/database/sql"
	model "ravxcheckout/src/internal/model/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func init() {

	database, err := sql.GetGormDB()
	if err != nil {
		return
	}
	db = database

	err = db.AutoMigrate(model.Order{})
	if err != nil {
		panic(err)
	}
}

type GetByIDFn func(ID string, notFull bool, payed bool) (*model.Order, error)

func GetByID(ID string, notFull bool, payed bool) (*model.Order, error) {
	result := &model.Order{}
	tx := db.Where("id = ?", ID)

	if !notFull {
		tx = tx.Preload("Customer.Address").Preload(clause.Associations)
	}
	if !payed {
		tx = tx.Where("status <> ?", "APPROVED")
	}

	tx.First(result)

	return result, tx.Error
}

type CreateFn func(user *model.Order) error

func Create(user *model.Order) error {
	tx := db.Create(user)
	return tx.Error
}

type UpdateFn func(user *model.Order) error

func Update(user *model.Order) error {
	tx := db.Save(user)
	return tx.Error
}

type PatchFn func(ID string, changes map[string]interface{}) error

func Patch(ID string, changes map[string]interface{}) error {
	delete(changes, "id")
	tx := db.Model(&model.Order{}).Where("id = ?", ID).Updates(changes)
	return tx.Error
}

type DeleteFn func(ID string) error

func Delete(ID string) error {
	tx := db.Delete(model.Order{}, "id = ?", ID)
	return tx.Error
}
