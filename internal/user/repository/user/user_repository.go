package user_repository

import (
	"context"
	"github.com/regiszanandrea/posty/internal/mongodb"
	"github.com/regiszanandrea/posty/internal/user/entity"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(user *entity.User) (string, error)
	Find(id string) (*entity.User, error)
	IncrementFollowers(id string) error
	IncrementFollowing(id string) error
	DecrementFollowers(id string) error
	DecrementFollowing(id string) error
	IncreasePostsCount(id string) error
}

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, configs *viper.Viper) *UserRepository {
	usersCollection := client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.user-collection"),
	)

	return &UserRepository{
		collection: usersCollection,
	}
}

func (repo *UserRepository) Create(user *entity.User) (string, error) {
	result, err := repo.collection.InsertOne(context.TODO(), user)

	if err != nil {
		if mongodb.IsDup(err) {
			return "", mongodb.ErrDuplicateKey
		}
		return "", err
	}

	objectId := result.InsertedID.(primitive.ObjectID)

	return objectId.Hex(), err
}

func (repo *UserRepository) Find(id string) (*entity.User, error) {
	var user entity.User

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	err = repo.collection.FindOne(
		context.TODO(), bson.M{"_id": objectId},
	).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) IncrementFollowers(id string) error {
	return repo.IncrementField(id, "followers_count", 1)
}

func (repo *UserRepository) IncrementFollowing(id string) error {
	return repo.IncrementField(id, "following_count", 1)
}

func (repo *UserRepository) DecrementFollowers(id string) error {
	return repo.IncrementField(id, "followers_count", -1)
}

func (repo *UserRepository) DecrementFollowing(id string) error {
	return repo.IncrementField(id, "following_count", -1)
}

func (repo *UserRepository) IncreasePostsCount(id string) error {
	return repo.IncrementField(id, "posts_count", 1)
}

func (repo *UserRepository) IncrementField(id, field string, value int) error {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = repo.collection.UpdateOne(context.TODO(), bson.M{"_id": objectId},
		bson.D{
			{"$inc",
				bson.D{
					{field, value},
				},
			},
		},
	)

	if err != nil {
		return err
	}

	return nil
}
