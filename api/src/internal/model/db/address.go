package modelDB

import "time"

type Address struct {
	ID         string `gorm:"primary_key;column:id;size:36;not null" json:"-"`
	Street     string `gorm:"column:street;size:60" json:"street"`
	Number     string `gorm:"column:number;size:10" json:"number"`
	Complement string `gorm:"column:complement;size:60" json:"complement"`
	District   string `gorm:"column:district;size:40" json:"district"`
	City       string `gorm:"column:city;size:40" json:"city"`
	State      string `gorm:"column:state;size:20" json:"state"`
	Country    string `gorm:"column:country;size:20" json:"country"`
	PostalCode string `gorm:"column:postal_code;size:20" json:"postal_code"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Address) TableName() string {
	return "addresses"
}
