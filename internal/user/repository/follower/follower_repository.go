package follower_repository

import (
	"context"
	"github.com/regiszanandrea/posty/internal/user/entity"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Follow(followerId, followingId string) (string, error)
	Unfollow(followerId, followingId string) error
	GetFollowingUsers(followerId string) ([]string, error)
}

type FollowerRepository struct {
	collection *mongo.Collection
}

func NewFollowerRepository(client *mongo.Client, configs *viper.Viper) *FollowerRepository {
	followersCollection := client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.follower-collection"),
	)

	return &FollowerRepository{
		collection: followersCollection,
	}
}

func (repo *FollowerRepository) Follow(followerId, followingId string) (string, error) {
	followerIdObjectId, err := primitive.ObjectIDFromHex(followerId)

	if err != nil {
		return "", err
	}

	followingIdObjectId, err := primitive.ObjectIDFromHex(followingId)

	if err != nil {
		return "", err
	}

	result, err := repo.collection.InsertOne(context.TODO(), bson.M{
		"follower_id": followerIdObjectId,
		"user_id":     followingIdObjectId,
	})

	if err != nil {
		return "", err
	}

	objectId := result.InsertedID.(primitive.ObjectID)

	return objectId.Hex(), err
}

func (repo *FollowerRepository) Unfollow(followerId, followingId string) error {
	followerIdObjectId, err := primitive.ObjectIDFromHex(followerId)

	if err != nil {
		return err
	}

	followingIdObjectId, err := primitive.ObjectIDFromHex(followingId)

	if err != nil {
		return err
	}

	_, err = repo.collection.DeleteOne(context.TODO(), bson.M{
		"follower_id": followerIdObjectId,
		"user_id":     followingIdObjectId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (repo *FollowerRepository) GetFollowingUsers(followerID string) ([]string, error) {
	objectId, err := primitive.ObjectIDFromHex(followerID)

	if err != nil {
		return nil, err
	}

	var result []string

	curr, err := repo.collection.Find(
		context.TODO(), bson.M{"follower_id": objectId},
		options.Find().SetProjection(bson.D{{"following_id", 0}}),
	)

	for curr.Next(context.TODO()) {
		var follower entity.Follower
		if err := curr.Decode(&follower); err != nil {
			return nil, err
		}

		result = append(result, follower.FollowingID.Hex())
	}

	return result, nil
}
