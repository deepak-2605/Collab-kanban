package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Board struct {
	ID        bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name      string          `bson:"name"          json:"name"`
	OwnerID   bson.ObjectID   `bson:"owner_id"      json:"ownerId"`
	Members   []bson.ObjectID `bson:"members"       json:"members"`
	CreatedAt time.Time       `bson:"created_at"    json:"createdAt"`
	UpdatedAt time.Time       `bson:"updated_at"    json:"updatedAt"`
}
