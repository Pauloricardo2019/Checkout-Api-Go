package user

import (
	"ravxcheckout/src/adapter/database/sql"
	"ravxcheckout/src/internal/model"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {

	database, err := sql.GetGormDB()
	if err != nil {
		return
	}
	db = database

	err = db.AutoMigrate(model.User{})
	if err != nil {
		panic(err)
	}
}

func GetByID(ID string) (*model.User, error) {
	result := &model.User{}
	tx := db.First(result, "id = ?", ID)

	return result, tx.Error
}

func Create(user *model.User) error {
	tx := db.Create(user)
	return tx.Error
}

func Update(user *model.User) error {
	tx := db.Save(user)
	return tx.Error
}

func Delete(ID string) error {
	tx := db.Delete(model.User{}, "id = ?", ID)
	return tx.Error
}
