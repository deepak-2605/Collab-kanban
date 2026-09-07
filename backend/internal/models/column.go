package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Column struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	BoardID   bson.ObjectID `bson:"board_id"      json:"boardId"`
	Name      string        `bson:"name"          json:"name"`
	Position  int           `bson:"position"      json:"position"`
	CreatedAt time.Time     `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at"    json:"updatedAt"`
}
