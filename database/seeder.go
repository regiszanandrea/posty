package main

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/regiszanandrea/posty/configs/app"
	"github.com/regiszanandrea/posty/internal/mongodb"
	post_entity "github.com/regiszanandrea/posty/internal/post/entity"
	"github.com/regiszanandrea/posty/internal/user/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"time"
)

var (
	postsCollection     *mongo.Collection
	usersCollection     *mongo.Collection
	followersCollection *mongo.Collection
)

func main() {
	configs := app.RegisterAppConfigs()

	client := mongodb.NewMongoDBClient(configs)

	err := client.Connect(context.Background())

	if err != nil {
		panic(err)
	}

	postsCollection = client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.post-collection"),
	)

	usersCollection = client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.user-collection"),
	)

	followersCollection = client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.follower-collection"),
	)

	users := createUsers(10)

	createFollowers(users)

	createPostForUsers(users, 5)
}

func createUsers(numberOfUsers int) []string {
	var users []interface{}
	var usersCreated []string

	for i := 0; i < numberOfUsers; i++ {
		user := entity.User{
			Username:       faker.Username(),
			CreatedAt:      time.Now(),
			FollowersCount: uint(rand.Uint32()),
			FollowingCount: uint(rand.Uint32()),
			PostsCount:     uint(rand.Uint32()),
		}

		users = append(users, user)
	}

	result, err := usersCollection.InsertMany(context.TODO(), users)

	if err != nil {
		panic(err)
	}

	for _, id := range result.InsertedIDs {
		usersCreated = append(usersCreated, id.(primitive.ObjectID).Hex())
	}

	return usersCreated
}

func createPostForUsers(users []string, numberOfPosts int) {

	for _, user := range users {
		objectId, _ := primitive.ObjectIDFromHex(user)
		var posts []interface{}

		for i := 0; i < numberOfPosts; i++ {

			post := post_entity.Post{
				UserID:    objectId,
				Content:   faker.Paragraph(),
				CreatedAt: time.Now(),
			}

			posts = append(posts, post)
		}

		_, err := postsCollection.InsertMany(context.TODO(), posts)

		if err != nil {
			panic(err)
		}
	}
}

func createFollowers(users []string) {
	max := len(users)

	for i := 0; i < len(users); i++ {
		followerId := users[rand.Intn(max-0)+0]
		followingId := users[rand.Intn(max-0)+0]

		if followingId == followerId {
			continue
		}

		followerIdObjectId, err := primitive.ObjectIDFromHex(followerId)

		if err != nil {
			panic(err)
		}

		followingIdObjectId, err := primitive.ObjectIDFromHex(followingId)

		if err != nil {
			panic(err)
		}

		_, err = followersCollection.InsertOne(context.TODO(), bson.M{
			"follower_id": followerIdObjectId,
			"user_id":     followingIdObjectId,
		})

		if err != nil {
			panic(err)
		}
	}
}
