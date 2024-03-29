package repository

import (
	"atommuse/backend/authentication-service/pkg/model"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(userID string) (*model.User, error)
	UpdateUserByID(userID string, updateUser *model.RequestUpdateUser) error
	UpdateUserPasswordByID(userID string, newPassword string) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) CreateUser(user *model.User) error {
	user.Role = "exhibitor"
	_, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(userID string) (*model.User, error) {
	// Convert the userID string to an ObjectId
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Define a filter to find the user by ID
	filter := bson.M{"_id": objectID}

	// Execute the find one operation
	var user model.User
	if err := r.collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if no document is found
		}
		return nil, err
	}

	return &user, nil
}

// updateUserByID updates a user by their ID.
func (r *userRepository) UpdateUserByID(userID string, updateUser *model.RequestUpdateUser) error {
	// Convert the string ID to ObjectId
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %v", err)
	}

	// Define the update fields
	update := bson.M{
		"$set": bson.M{
			"username":  updateUser.UserName,
			"firstname": updateUser.FirstName,
			"lastname":  updateUser.LastName,
			"email":     updateUser.Email,
			"profile":   updateUser.ProfileImage,
		},
	}

	// Prepare the filter for the update
	filter := bson.M{"_id": objectID}

	// Perform the update operation
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *userRepository) UpdateUserPasswordByID(userID string, newPassword string) error {

	// Convert the string ID to ObjectId
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %v", err)
	}
	// Create an update to set the new password
	update := bson.M{"$set": bson.M{"password": newPassword}}

	// Prepare the filter for the update
	filter := bson.M{"_id": objectID}

	// Perform the update operation
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return err
}
