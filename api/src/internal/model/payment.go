package model

import (
	modelDB "ravxcheckout/src/internal/model/db"
	"time"
)

type PostData struct {
	IssuerPaymentID            string `json:"issuer_payment_id"`
	PayerAuthenticationRequest string `json:"payer_authentication_request"`
}

type DebitResponse struct {
	AuthorizationCode     string    `json:"authorization_code"`
	AuthorizedTimestamp   time.Time `json:"authorized_timestamp"`
	ReasonCode            string    `json:"reason_code"`
	ReasonMessage         string    `json:"reason_message"`
	Acquirer              string    `json:"acquirer"`
	SoftDescriptor        string    `json:"soft_descriptor"`
	Brand                 string    `json:"brand"`
	TerminalNsu           string    `json:"terminal_nsu"`
	AcquirerTransactionID string    `json:"acquirer_transaction_id"`
	TransactionID         string    `json:"transaction_id"`
}

type Card struct {
	Number          string `json:"number"`
	CardHolderName  string `json:"cardholder_name,omitempty"`
	SecurityCode    string `json:"security_code,omitempty"`
	Brand           string `json:"brand,omitempty"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
}

type Debit struct {
	CardHolderMobile string `json:"cardholder_mobile,omitempty"`
	SoftDescription  string `json:"soft_descriptor,omitempty"`
	DynamicMcc       int32  `json:"dynamic_mcc,omitempty"`
	Authenticated    bool   `json:"authenticated"`
	Card             *Card  `json:"card"`
}

type Order struct {
	OrderID     string  `json:"order_id"`
	Sales_tax   float32 `json:"sales_tax,omitempty"`
	ProductType string  `json:"product_type,omitempty"`
}

type Address struct {
	Street     string `json:"street,omitempty"`
	Number     string `json:"number,omitempty"`
	Complement string `json:"complement,omitempty"`
	District   string `json:"district,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

type Customer struct {
	CustomerID     string   `json:"customer_id"`
	FirstName      string   `json:"first_name,omitempty"`
	LastName       string   `json:"last_name,omitempty"`
	Name           string   `json:"name,omitempty"`
	Email          string   `json:"email,omitempty"`
	DocumentType   string   `json:"document_type,omitempty"`
	DocumentNumber string   `json:"document_number,omitempty"`
	PhoneNumber    string   `json:"phone_number,omitempty"`
	Address        *Address `json:"billing_address"`
}

type Device struct {
	IpAddress string `json:"ip_address,omitempty"`
}

type Shippings struct {
	FirstName      string   `json:"first_name,omitempty"`
	Name           string   `json:"name,omitempty"`
	Email          string   `json:"email,omitempty"`
	PhoneNumber    string   `json:"phone_number,omitempty"`
	ShippingAmount uint32   `json:"shipping_amount,omitempty"`
	Address        *Address `json:"address,omitempty"`
}

type SubMerchant struct {
	IdentificationCode string `json:"identification_code,omitempty"`
	DocumentType       string `json:"document_type,omitempty"`
	DocumentNumber     string `json:"document_number,omitempty"`
	Address            string `json:"address,omitempty"`
	City               string `json:"city,omitempty"`
	State              string `json:"state,omitempty"`
	PostalCode         string `json:"postal_code,omitempty"`
}
type PaymentResponse struct {
	PaymentID   string         `json:"payment_id,omitempty"`
	SellerID    string         `json:"seller_id,omitempty"`
	RedirectUrl string         `json:"redirect_url,omitempty"`
	PostData    *PostData      `json:"post_data,omitempty"`
	Amount      uint           `json:"amount,omitempty"`
	Currency    string         `json:"currency,omitempty"`
	OrderID     string         `json:"order_id,omitempty"`
	Status      string         `json:"status,omitempty"`
	ReceivedAt  *time.Time     `json:"received_at,omitempty"`
	Debit       *DebitResponse `json:"debit,omitempty"`
	Error       *Error         `json:"error,omitempty"`
}

type PaymentRequest struct {
	OrderID string         `json:"order_id"`
	Order   *modelDB.Order `json:"order,omitempty"`
	Debit   *Debit         `json:"debit,required"`
}

type PixPaymentRequest struct {
	OrderID string         `json:"order_id"`
	Order   *modelDB.Order `json:"order,omitempty"`
}

type PixPaymentResponse struct {
	PaymentID      string             `json:"payment_id,omitempty"`
	Status         string             `json:"status,omitempty"`
	Description    string             `json:"description,omitempty"`
	AdditionalData *AdditionalDataPix `json:"additional_data,omitempty"`
	Error          *Error             `json:"error,omitempty"`
}

type AdditionalDataPix struct {
	TransactionID        string `json:"transaction_id"`
	QrCode               string `json:"qr_code"`
	CreationDateQrCode   string `json:"creation_date_qrcode"`
	ExpirationDateQrCode string `json:"expiration_date_qrcode"`
	PSPCode              string `json:"psp_code"`
}
