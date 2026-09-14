package repository

import (
	"context"

	"github.com/deepak-2605/collab-kanban/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CardRepository struct {
	coll *mongo.Collection
}

func NewCardRepository(db *mongo.Database) *CardRepository {
	return &CardRepository{coll: db.Collection("cards")}
}

// Create — InsertOne, return err.
func (r *CardRepository) Create(ctx context.Context, card *models.Card) error {

	res, err := r.coll.InsertOne(ctx, card)

	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		card.ID = oid
	}
	return nil
}

// FindByID —  FindOne with bson.M{"_id": id}.
func (r *CardRepository) FindByID(ctx context.Context, id bson.ObjectID) (*models.Card, error) {

	var card models.Card

	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&card)

	if err != nil {
		return nil, err
	}
	return &card, nil
}

// FindByColumn — Returns all cards for a column, ordered by position ascending.
func (r *CardRepository) FindByColumn(ctx context.Context, columnID bson.ObjectID) ([]models.Card, error) {

	opts := options.Find().SetSort(bson.M{"position": 1})

	cursor, err := r.coll.Find(ctx, bson.M{"column_id": columnID}, opts)
	if err != nil {
		return nil, err
	}

	cards := []models.Card{}

	if err := cursor.All(ctx, &cards); err != nil {
		return nil, err
	}

	return cards, nil
}

// Update - update the mutable fields of an existing card.
func (r *CardRepository) Update(ctx context.Context, card *models.Card) error {
	_, err := r.coll.UpdateOne(ctx,
		bson.M{"_id": card.ID},
		bson.M{"$set": bson.M{
			"title":       card.Title,
			"description": card.Description,
			"column_id":   card.ColumnID, // ← changes when a card MOVES to another column
			"position":    card.Position, // ← changes on move/reorder
			"assignee_id": card.AssigneeID,
			"updated_at":  card.UpdatedAt,
		}})

	return err
}

// Delete - Card for a particular column
func (r *CardRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
