package modelDB

import "time"

type Customer struct {
	ID             string   `gorm:"primary_key;column:id;size:36" json:"customer_id"`
	FirstName      string   `gorm:"column:first_name;size:40" json:"first_name"`
	LastName       string   `gorm:"column:last_name;size:80" json:"last_name"`
	Name           string   `gorm:"column:name;size:100" json:"name"`
	Email          string   `gorm:"column:email;size:255" json:"email"`
	DocumentType   string   `gorm:"column:document_type;size:26" json:"document_type"`
	DocumentNumber string   `gorm:"column:document_number;size:25" json:"document_number"`
	PhoneNumber    string   `gorm:"column:phone_number;size:25" json:"phone_number"`
	AddressID      string   `gorm:"column:address_id;size:36" json:"-"`
	Address        *Address `gorm:"foreign_key:address_id,references:id" json:"billing_address"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Customer) TableName() string {
	return "customers"
}
