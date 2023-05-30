package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type Source struct {
	mgm.DefaultModel `bson:",inline"`
	ID               *string   `json:"_id" bson:"_id,omitempty"`
	Name             *string   `json:"name" bson:"name"`
	Url              *string   `json:"url" bson:"url"`
	Category         []*string `json:"category" bson:"category"`
	Icon             *string   `json:"icon" bson:"icon"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}
