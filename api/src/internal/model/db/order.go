package modelDB

import "time"

type Order struct {
	ID          string    `gorm:"primary_key;column:id;size:36" json:"order_id"`
	Status      string    `gorm:"column:status;size:25" json:"status"`
	ProductType string    `gorm:"column:status;size:20" json:"product_type"`
	Amount      uint      `gorm:"column:amount;not null" json:"amount"`
	Currency    string    `gorm:"column:currency;size:3;not null" json:"currency"`
	CustomerID  string    `gorm:"column:customer_id;size:36" json:"-"`
	Customer    *Customer `gorm:"foreign_key:customer_id,references:id" json:"customer"`
	PaymentID   *string   `gorm:"column:payment_id;size:36" json:"-"`
	RedirectUrl *string   `gorm:"-" json:"redirect_url,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Order) TableName() string {
	return "orders"
}
