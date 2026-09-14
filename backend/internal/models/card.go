package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Card struct {
	ID          bson.ObjectID  `bson:"_id,omitempty"          json:"id"`
	BoardID     bson.ObjectID  `bson:"board_id"              json:"boardId"`
	ColumnID    bson.ObjectID  `bson:"column_id"             json:"columnId"`
	Title       string         `bson:"title"                 json:"title"`
	Description string         `bson:"description"           json:"description"`
	Position    int            `bson:"position"              json:"position"`
	AssigneeID  *bson.ObjectID `bson:"assignee_id,omitempty" json:"assigneeId,omitempty"`
	CreatedAt   time.Time      `bson:"created_at"            json:"createdAt"`
	UpdatedAt   time.Time      `bson:"updated_at"            json:"updatedAt"`
}
