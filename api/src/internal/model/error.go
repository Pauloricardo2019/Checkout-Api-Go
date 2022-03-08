package model

import "time"

type Error struct {
	Message    string    `json:"message"`
	Name       string    `json:"name"`
	StatusCode int       `json:"status_code"`
	Details    []Details `json:"details"`
}

type ErrorReduced struct {
	Name        string `json:"error"`
	Description string `json:"error_description"`
}

type Details struct {
	Status            string     `json:"status,omitempty"`
	ErrorCode         string     `json:"error_code,omitempty"`
	Description       string     `json:"description,omitempty"`
	DescriptionDetail string     `json:"description_detail,omitempty"`
	Antifraud         *Antifraud `json:"antifraud,omitempty"`
}

type Antifraud struct {
	Code                     string    `json:"code"`
	Description              string    `json:"description"`
	TransactionToken         string    `json:"transaction_token"`
	TransactionID            string    `json:"transaction_id"`
	TransactionReferenceCode string    `json:"transaction_reference_code"`
	TransactionDatetime      time.Time `json:"transaction_datetime"`
}
