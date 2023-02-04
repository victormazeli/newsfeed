package models

import "github.com/kamva/mgm/v3"

type News struct {
	mgm.DefaultModel `bson:",inline"`
	ID               string `json:"_id" bson:"_id,omitempty"`
	Author           string `json:"author"`
	title            string `json:"title"`
	description      string `json:"description"`
	url              string `json:"url"`
	urlToImage       string `json:"url_to_image"`
	PublishedAt      bool   `json:"published_at"`
	Content          string `json:"content"`
}
