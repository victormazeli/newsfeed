package models

import "github.com/kamva/mgm/v3"

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string `json:"email"`
	Picture          string `json:"picture"`
	FullName         string `json:"full_name"`
	UserId           string `json:"user_id"`
	Topics           []Topic
	UpdatedAt        string `json:"updated_at" bson:"updated_at"`
	CreatedAt        string `json:"created_at" bson:"created_at"`
}

type Topic struct {
	Topic string `json:"topic"`
}
