package repository

import (
	"context"

	"github.com/deepak-2605/collab-kanban/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BoardRepository struct {
	coll *mongo.Collection
}

func NewBoardRepository(db *mongo.Database) *BoardRepository {
	return &BoardRepository{coll: db.Collection("boards")}
}

// FindByUser — returns every board where the user is EITHER the owner OR in the members list.
func (r *BoardRepository) FindByUser(ctx context.Context, userID bson.ObjectID) ([]models.Board, error) {
	filter := bson.M{"$or": []bson.M{
		{"owner_id": userID},
		{"members": userID},
	}}

	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var boards []models.Board
	if err := cursor.All(ctx, &boards); err != nil {
		return nil, err
	}
	return boards, nil
}

// Create — InsertOne, return err.
func (r *BoardRepository) Create(ctx context.Context, board *models.Board) error {

	res, err := r.coll.InsertOne(ctx, board)

	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		board.ID = oid
	}
	return nil
}

// FindByID —  FindOne with bson.M{"_id": id}.
func (r *BoardRepository) FindByID(ctx context.Context, id bson.ObjectID) (*models.Board, error) {

	var board models.Board

	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&board)

	if err != nil {
		return nil, err
	}
	return &board, nil
}

// Update - update the mutable fields of an existing board.
func (r *BoardRepository) Update(ctx context.Context, board *models.Board) error {
	_, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": board.ID},
		bson.M{"$set": bson.M{
			"name":       board.Name,
			"members":    board.Members,
			"updated_at": board.UpdatedAt,
		}})

	return err
}

// Delete - BSON Record
func (r *BoardRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
