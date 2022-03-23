package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	. "go.uber.org/fx"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	usersCollectionIndexes = []mongo.IndexModel{
		{
			Keys: bson.M{
				"username": -1,
			}, Options: options.Index().SetUnique(true),
		},
	}
	postsCollectionIndexes = []mongo.IndexModel{
		{
			Keys: bson.M{
				"user_id": -1,
			}, Options: nil,
		},
		{
			Keys: bson.D{
				{"user_id", -1},
				{"created_at", -1},
			}, Options: nil,
		},
	}
	followersCollectionIndexes = []mongo.IndexModel{
		{
			Keys: bson.M{
				"follower_id": -1,
			}, Options: nil,
		},
		{
			Keys: bson.M{
				"user_id": -1,
			}, Options: nil,
		},
	}
	ErrDuplicateKey = errors.New("there is already a key with this value")
)

func NewMongoDBClient(configs *viper.Viper) *mongo.Client {
	connectUrl := "mongodb://" +
		configs.GetString("app.mongodb.user") + ":" +
		configs.GetString("app.mongodb.password") + "@" +
		configs.GetString("app.mongodb.host") + ":" +
		configs.GetString("app.mongodb.port")

	clientOptions := options.Client().ApplyURI(connectUrl)

	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		panic(err)
	}

	return client
}

func RegisterMongoDB(lifecycle Lifecycle, client *mongo.Client, configs *viper.Viper) {
	lifecycle.Append(Hook{
		OnStart: func(ctx context.Context) error {

			err := client.Connect(ctx)

			if err != nil {
				return err
			}

			err = CreateIndexes(client, configs, ctx)

			if err != nil {
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})
}

func CreateIndexes(client *mongo.Client, configs *viper.Viper, ctx context.Context) error {
	err := CreateUsersCollectionIndexes(client, configs, ctx)

	if err != nil {
		return err
	}

	err = CreateFollowersCollectionIndexes(client, configs, ctx)

	if err != nil {
		return err
	}

	return CreatePostsCollectionIndexes(client, configs, ctx)
}

func CreateUsersCollectionIndexes(client *mongo.Client, configs *viper.Viper, ctx context.Context) error {
	usersCollection := client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.user-collection"),
	)

	err := CreateCollectionIndexes(usersCollection, usersCollectionIndexes, ctx)
	if err != nil {
		return err
	}

	return nil
}

func CreateFollowersCollectionIndexes(client *mongo.Client, configs *viper.Viper, ctx context.Context) error {
	followersCollection := client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.follower-collection"),
	)

	err := CreateCollectionIndexes(followersCollection, followersCollectionIndexes, ctx)
	if err != nil {
		return err
	}

	return nil
}

func CreatePostsCollectionIndexes(client *mongo.Client, configs *viper.Viper, ctx context.Context) error {
	postsCollection := client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.post-collection"),
	)

	err := CreateCollectionIndexes(postsCollection, postsCollectionIndexes, ctx)
	if err != nil {
		return err
	}

	return nil
}

func CreateCollectionIndexes(collection *mongo.Collection, collectionIndexes []mongo.IndexModel, ctx context.Context) error {
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return err
	}

	var indices []mongo.IndexSpecification
	err = cursor.All(nil, &indices)

	if err != nil {
		return err
	}

	// this +1 references to _id index that is already created when you create the collection
	if len(indices) == len(collectionIndexes)+1 {
		return nil
	}

	_, err = collection.Indexes().CreateMany(ctx, collectionIndexes)

	if err != nil {
		return err
	}

	return nil
}

func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}
