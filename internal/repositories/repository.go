package repositories

import (
	"context"
	"go-quickstart/internal/canonical"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	CreateUser(user canonical.User) error
	GetAllUsers() ([]canonical.User, error)
	GetUserById(id string) (canonical.User, error)
	UpdateUser(id string, user canonical.User) error
	DeleteUser(id string) error
}

type repository struct {
	collection *mongo.Collection
}

func New() Repository {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:password@localhost:27017"))
	if err != nil {
		panic(err)
	}

	return &repository{
		collection: client.Database("user-database").Collection("userSlice"),
	}
}

func (repo *repository) CreateUser(user canonical.User) error {
	_, err := repo.collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) GetAllUsers() ([]canonical.User, error) {
	var userSlice []canonical.User

	res, err := repo.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	for res.Next(context.Background()) {
		var user canonical.User

		err := res.Decode(&user)
		if err != nil {
			return nil, err
		}

		userSlice = append(userSlice, user)
	}

	return userSlice, nil
}

func (repo *repository) GetUserById(id string) (canonical.User, error) {
	var user canonical.User

	err := repo.collection.FindOne(context.Background(), bson.D{
		{
			Key:   "_id",
			Value: id,
		},
	}).Decode(&user)

	if err != nil {
		return canonical.User{}, err
	}

	return user, nil
}

func (repo *repository) UpdateUser(id string, user canonical.User) error {
	filter := bson.D{{Key: "_id", Value: id}}
	fields := bson.M{
		"$set": bson.M{
			"name": user.Name,
		},
	}

	_, err := repo.collection.UpdateOne(context.Background(), filter, fields)

	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) DeleteUser(id string) error {
	_, err := repo.collection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}

	return nil
}
