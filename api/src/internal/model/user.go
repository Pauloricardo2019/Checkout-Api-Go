package model

type User struct {
	ID    string `json:"_id" bson:"_id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
