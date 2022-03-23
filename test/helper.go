package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/regiszanandrea/posty/internal"
	"github.com/regiszanandrea/posty/internal/mongodb"
	"github.com/regiszanandrea/posty/internal/post"
	post_entity "github.com/regiszanandrea/posty/internal/post/entity"
	"github.com/regiszanandrea/posty/internal/user"
	"github.com/regiszanandrea/posty/internal/user/entity"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"net/http"
	"net/url"
	"time"
)

var (
	MongoClient *mongo.Client
)

func SetUpFXApp() *fx.App {
	return fx.New(
		fx.NopLogger,
		internal.ApplicationModule,
		fx.Invoke(RegisterMongoDB),
		user.Invokables,
		post.Invokables,
		fx.Invoke(func(fa *fiber.App, c *viper.Viper) {
			go Boot(c)(fa)
		}),
	)
}

func Boot(c *viper.Viper) func(app *fiber.App) error {
	return func(app *fiber.App) error {
		if err := app.Listen(c.GetString("app.fiber.address")); err != nil {
			return err
		}

		return nil
	}
}

func RegisterMongoDB(client *mongo.Client, configs *viper.Viper) {
	err := client.Connect(context.Background())

	if err != nil {
		panic(err)
	}

	err = mongodb.CreateIndexes(client, configs, context.Background())

	if err != nil {
		panic(err)
	}

	MongoClient = client
}

func StopFXApp(app *fx.App) {
	appStopCtx, appStopCancel := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)

	if err := app.Stop(appStopCtx); err != nil {
		panic(err)
	}

	appStopCancel()
}

func MakePostRequest(host string, endpoint string, body map[string]string) *http.Response {
	postBody, _ := json.Marshal(body)
	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(
		"http://"+host+endpoint,
		"application/json",
		requestBody,
	)

	if err != nil {
		panic(err)
	}

	return resp
}

func MakeGetRequest(host string, endpoint string, queryParams map[string]string) *http.Response {
	base, _ := url.Parse("http://" + host + endpoint)

	params := url.Values{}

	for key, val := range queryParams {
		params.Add(key, val)
	}

	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())

	if err != nil {
		panic(err)
	}

	return resp
}

func CreateUsers(numberOfUsers int, usersCollection *mongo.Collection) *mongo.InsertManyResult {
	var users []interface{}

	for i := 0; i < numberOfUsers; i++ {
		u := entity.User{
			Username:  faker.Username(),
			CreatedAt: time.Now(),
		}

		users = append(users, u)
	}

	results, _ := usersCollection.InsertMany(context.Background(), users)
	return results
}

func CreatePosts(numberOfPosts int, userId primitive.ObjectID, postsCollection *mongo.Collection) *mongo.InsertManyResult {
	var users []interface{}

	for i := 0; i < numberOfPosts; i++ {
		u := post_entity.Post{
			UserID:    userId,
			Content:   faker.Paragraph(),
			CreatedAt: time.Now(),
		}

		users = append(users, u)
	}

	results, _ := postsCollection.InsertMany(context.Background(), users)
	return results
}

func CreateFollower(followerId primitive.ObjectID, followingId primitive.ObjectID, followersCollection *mongo.Collection) *mongo.InsertOneResult {
	follower := entity.Follower{
		FollowerID:  followerId,
		FollowingID: followingId,
		CreatedAt:   time.Now(),
	}

	result, _ := followersCollection.InsertOne(context.Background(), follower)

	return result
}
