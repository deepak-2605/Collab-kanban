package repository

import (
	"context"

	"github.com/deepak-2605/collab-kanban/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ColumnRepository struct {
	coll *mongo.Collection
}

func NewColumnRepository(db *mongo.Database) *ColumnRepository {
	return &ColumnRepository{coll: db.Collection("boards")}
}

// Create — InsertOne, return err.
func (r *ColumnRepository) Create(ctx context.Context, column *models.Column) error {

	res, err := r.coll.InsertOne(ctx, column)

	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		column.ID = oid
	}
	return nil
}

// FindByID —  FindOne with bson.M{"_id": id}.
func (r *ColumnRepository) FindByID(ctx context.Context, id bson.ObjectID) (*models.Column, error) {

	var column models.Column

	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&column)

	if err != nil {
		return nil, err
	}
	return &column, nil
}

// Update - update the mutable fields of an existing column.
func (r *ColumnRepository) Update(ctx context.Context, column *models.Column) error {
	_, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": column.ID},
		bson.M{"$set": bson.M{
			"name":       column.Name,
			"position":   column.Position,
			"updated_at": column.UpdatedAt,
		}})

	return err
}

// Delete - BSON Record
func (r *ColumnRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// FindByBoard — Returns all columns for a board, ordered by position ascending.
func (r *ColumnRepository) FindByBoard(ctx context.Context, boardID bson.ObjectID) ([]models.Column, error) {

	opts := options.Find().SetSort(bson.M{"position": 1})

	cursor, err := r.coll.Find(ctx, bson.M{"board_id": boardID}, opts)
	if err != nil {
		return nil, err
	}

	columns := []models.Column{}

	if err := cursor.All(ctx, &columns); err != nil {
		return nil, err
	}

	return columns, nil

}
