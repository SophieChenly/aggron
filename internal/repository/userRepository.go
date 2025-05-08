package repository

import (
	"aggron/internal/db/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, discordID, email string) (*models.User, error) {
	existingUser, _ := r.FindByDiscordID(ctx, discordID)
	if existingUser != nil {
		return nil, errors.New("user with this Discord ID and email already exists")
	}

	user := models.User{
		DiscordID: discordID,
		Email:     email,
	}

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateEmail(ctx context.Context, discordID, newEmail string) (*models.User, error) {
	_, err := r.FindByDiscordID(ctx, discordID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"discord_id": discordID}
	update := bson.M{
		"$set": bson.M{
			"email": newEmail,
		},
	}

	// get back updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser models.User

	err = r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *UserRepository) FindByDiscordID(ctx context.Context, discordID string) (*models.User, error) {
	var user models.User
	filter := bson.M{"discord_id": discordID}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
