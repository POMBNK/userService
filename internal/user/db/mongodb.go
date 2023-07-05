package db

import (
	"context"
	"fmt"
	"github.com/POMBNK/restAPI/internal/pkg/apierror"
	"github.com/POMBNK/restAPI/internal/user"
	"github.com/POMBNK/restAPI/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	idKey = "_id"
)

type db struct {
	collection *mongo.Collection
	logs       *logger.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logs.Debug("Create user...")
	res, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("error, failed to create user:%w", err)
	}
	d.logs.Debug("Convert created USER ID to object ID...")
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("error, failed to cast result to object id")
	}
	return oid.Hex(), nil
}

func (d *db) GetById(ctx context.Context, id string) (user.User, error) {
	var res user.User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return res, fmt.Errorf("error failed to convert ID hex to object ID")
	}
	filter := bson.M{idKey: oid}
	err = d.collection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			d.logs.Debug("No user found")
			return res, apierror.ErrNotFound
		}
		return res, fmt.Errorf("failed to found user with id:%s due error: %w", id, err)
	}
	return res, nil
}

func (d *db) GetAll(ctx context.Context) ([]user.User, error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("can't find fields due with error: %w ", err)
	}
	var listOfUsers []user.User

	if err = cursor.All(ctx, &listOfUsers); err != nil {
		return nil, fmt.Errorf("can't get list of fields error: %w", err)
	}
	return listOfUsers, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("error failed to convert ID hex to object ID")
	}

	filter := bson.M{idKey: oid}
	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user to bson error: %w", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user obj to bson error: %w", err)
	}
	delete(updateUserObj, idKey) // cause _id autogenerating protected field

	update := bson.M{"$set": updateUserObj}
	res, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user obj error: %w", err)
	}

	if res.MatchedCount == 0 {
		return apierror.ErrNotFound
	}

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error failed to convert ID hex to object ID")
	}
	filter := bson.M{idKey: oid}
	res, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user error:%w", err)
	}

	if res.DeletedCount == 0 {
		return apierror.ErrNotFound
	}

	return nil
}

func New(database *mongo.Database, collection string, logs *logger.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logs:       logs,
	}
}
