package dto

type OrderDTO struct {
	RedirectUrl string `json:"redirect_url,omitempty"`
	Status      string `json:"status,omitempty"`
}
