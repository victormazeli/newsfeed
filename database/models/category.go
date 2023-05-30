package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type Category struct {
	mgm.DefaultModel `bson:",inline"`
	ID               *string   `json:"_id" bson:"_id,omitempty"`
	Name             *string   `json:"name" bson:"name"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}
