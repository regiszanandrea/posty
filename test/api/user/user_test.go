package user_test

import (
	"context"
	"encoding/json"
	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/regiszanandrea/posty/configs/app"
	"github.com/regiszanandrea/posty/internal/user/entity"
	"github.com/regiszanandrea/posty/test"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"testing"
)

var (
	application         *fx.App
	configs             *viper.Viper
	followersCollection *mongo.Collection
	usersCollection     *mongo.Collection
)

func TestUserApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User API test suite")
}

var _ = BeforeSuite(func() {
	application = helper.SetUpFXApp()

	configs = app.RegisterAppConfigs()

	usersCollection = helper.MongoClient.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.user-collection"),
	)

	followersCollection = helper.MongoClient.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.follower-collection"),
	)
})

var _ = AfterSuite(func() {
	helper.StopFXApp(application)

	_, err := usersCollection.DeleteMany(context.Background(), bson.M{})

	if err != nil {
		panic(err)
	}

	_, err = followersCollection.DeleteMany(context.Background(), bson.M{})

	if err != nil {
		panic(err)
	}
})

var _ = Describe("User API test", func() {
	Describe("Creating user", func() {
		Context("when its given a valid user", func() {
			It("creates it without error", func() {
				requestBody := map[string]string{
					"username": faker.Username(),
				}

				resp := helper.MakePostRequest(configs.GetString("app.fiber.address"), "/users", requestBody)

				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusCreated))
			})
		})

		Context("when already is a user with same username", func() {
			It("returns error", func() {
				requestBody := map[string]string{
					"username": faker.Username(),
				}

				helper.MakePostRequest(configs.GetString("app.fiber.address"), "/users", requestBody)
				resp := helper.MakePostRequest(configs.GetString("app.fiber.address"), "/users", requestBody)
				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusBadRequest))
			})
		})
	})

	Describe("Finding a user", func() {
		Context("when its given a valid user", func() {
			It("returns it", func() {
				requestBody := map[string]string{
					"username": faker.Username(),
				}
				var user, userReturned entity.User

				resp := helper.MakePostRequest(configs.GetString("app.fiber.address"), "/users", requestBody)

				json.NewDecoder(resp.Body).Decode(&user)

				resp = helper.MakeGetRequest(configs.GetString("app.fiber.address"), "/users/"+user.ID.Hex(), nil)

				json.NewDecoder(resp.Body).Decode(&userReturned)

				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusOK))
				Expect(user.ID).To(Equal(userReturned.ID))
			})
		})

		Context("when its given a non-existent user", func() {
			It("returns not found", func() {
				resp := helper.MakeGetRequest(configs.GetString("app.fiber.address"), "/users/"+primitive.NewObjectID().Hex(), nil)

				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusNotFound))
			})
		})
	})
	Describe("Following a user", func() {
		Context("when its given a follower and the followed user", func() {
			It("follows without error", func() {
				results := helper.CreateUsers(2, usersCollection)

				followerId := results.InsertedIDs[0].(primitive.ObjectID).Hex()
				followingId := results.InsertedIDs[1].(primitive.ObjectID).Hex()

				endpoint := "/users/" + followerId + "/follow/" + followingId

				resp := helper.MakePostRequest(configs.GetString("app.fiber.address"), endpoint, nil)

				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusOK))

				assertNumberOfFollowing(followerId, 1)
				assertNumberOfFollowers(followingId, 1)
			})

		})

		Context("when its given same user as follower and following", func() {
			It("returns error", func() {
				userId := primitive.NewObjectID().Hex()

				endpoint := "/users/" + userId + "/follow/" + userId

				resp := helper.MakePostRequest(configs.GetString("app.fiber.address"), endpoint, nil)
				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusBadRequest))
			})
		})
	})

	Describe("Unfollowing a user", func() {
		Context("when its given a follower and the followed user", func() {
			It("unfollows without error", func() {
				results := helper.CreateUsers(2, usersCollection)

				followerId := results.InsertedIDs[0].(primitive.ObjectID).Hex()
				followingId := results.InsertedIDs[1].(primitive.ObjectID).Hex()

				endpoint := "/users/" + followerId + "/follow/" + followingId

				helper.MakePostRequest(configs.GetString("app.fiber.address"), endpoint, nil)

				endpoint = "/users/" + followerId + "/unfollow/" + followingId

				resp := helper.MakePostRequest(configs.GetString("app.fiber.address"), endpoint, nil)
				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusOK))

				assertNumberOfFollowing(followerId, 0)
				assertNumberOfFollowers(followingId, 0)
			})
		})
	})
})

func assertNumberOfFollowers(followingId string, expected int) {
	resp := helper.MakeGetRequest(configs.GetString("app.fiber.address"), "/users/"+followingId, nil)

	var userFollowingReturned entity.User
	json.NewDecoder(resp.Body).Decode(&userFollowingReturned)

	Expect(userFollowingReturned.FollowersCount).To(BeEquivalentTo(expected))
}

func assertNumberOfFollowing(followerId string, expected int) {
	resp := helper.MakeGetRequest(configs.GetString("app.fiber.address"), "/users/"+followerId, nil)

	var userFollowerReturned entity.User
	json.NewDecoder(resp.Body).Decode(&userFollowerReturned)

	Expect(userFollowerReturned.FollowingCount).To(BeEquivalentTo(expected))
}
