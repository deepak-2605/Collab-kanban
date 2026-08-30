package repository

import (
	"context"

	"github.com/deepak-2605/collab-kanban/backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection("users")}
}

// FindByEmail
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create — It creates a user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {

	_, err := r.coll.InsertOne(ctx, user)

	return err
}

// FindByID
func (r *UserRepository) FindByID(ctx context.Context, id bson.ObjectID) (*models.User, error) {

	var user models.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}
