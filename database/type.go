package database

import (
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type UTCTime struct {
	time.Time
}

func (t UTCTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	utcTime := t.Time.UTC()
	return bson.MarshalValue(utcTime)
}
