package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type News struct {
	mgm.DefaultModel `bson:",inline"`
	ID               *string   `json:"_id" bson:"_id,omitempty"`
	Creator          []*string `json:"creator" bson:"creator"`
	Category         []*string `json:"category" bson:"category"`
	Title            *string   `json:"title" bson:"title"`
	Description      *string   `json:"description" bson:"description"`
	ImageURL         *string   `json:"image_url" bson:"image_url"`
	Link             *string   `json:"link" bson:"link"`
	PubDate          *string   `json:"pub_date" bson:"pub_date"`
	Content          *string   `json:"content" bson:"content"`
	SourceID         *string   `json:"source_id" bson:"source_id"`
	Likes            []*string `json:"likes" bson:"likes,default:[]"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}
