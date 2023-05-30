package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SavedNews struct {
	mgm.DefaultModel `bson:",inline"`
	NewsID           primitive.ObjectID `json:"news_id" bson:"news_id"`
	UserID           primitive.ObjectID `json:"user_id" bson:"user_id"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}
