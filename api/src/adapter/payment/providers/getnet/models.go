package getnet

import (
	"ravxcheckout/src/internal/model"
)

type CardTokenization struct {
	CardNumber string `json:"card_number"`
}

type Card struct {
	NumberToken     string `json:"number_token"`
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

type Device struct {
	IpAddress string `json:"ip_address,omitempty"`
	DeviceID  string `json:"device_id,omitempty"`
}

type Transaction struct {
	SellerID    string             `json:"seller_id"`
	Amount      uint               `json:"amount"`
	Currency    string             `json:"currency,omitempty"`
	Order       *model.Order       `json:"order"`
	Customer    *model.Customer    `json:"customer"`
	Device      *Device            `json:"device,omitempty"`
	Shippings   []model.Shippings  `json:"shippings,omitempty"`
	SubMerchant *model.SubMerchant `json:"sub_merchant,omitempty"`
	Debit       *Debit             `json:"debit"`
}

type AuthResponse struct {
	AccessToken string              `json:"access_token"`
	TokenType   string              `json:"token_type"`
	ExpiresIn   int                 `json:"expires_in"`
	Scope       string              `json:"scope"`
	Error       *model.ErrorReduced `json:"error,omitempty"`
}

type TokenizeCardResponse struct {
	NumberToken string       `json:"number_token"`
	Error       *model.Error `json:"error,omitempty"`
}

type PixTransaction struct {
	Amount     uint   `json:"amount"`
	Currency   string `json:"currency"`
	OrderID    string `json:"order_id,omitempty"`
	CustomerId string `json:"customer_id,omitempty"`
}
