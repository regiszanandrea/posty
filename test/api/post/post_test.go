package post_test

import (
	"context"
	"encoding/json"
	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/regiszanandrea/posty/configs/app"
	"github.com/regiszanandrea/posty/internal/post/entity"
	helper "github.com/regiszanandrea/posty/test"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	application         *fx.App
	configs             *viper.Viper
	postsCollection     *mongo.Collection
	followersCollection *mongo.Collection
	usersCollection     *mongo.Collection
)

func TestPostApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Post API test suite")
}

var _ = BeforeSuite(func() {
	application = helper.SetUpFXApp()

	configs = app.RegisterAppConfigs()

	postsCollection = helper.MongoClient.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.post-collection"),
	)

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

	_, err = postsCollection.DeleteMany(context.Background(), bson.M{})

	if err != nil {
		panic(err)
	}
})

var _ = Describe("Post API test", func() {
	Describe("Creating post", func() {
		Context("when its given a valid post", func() {
			It("creates it without error", func() {
				user := helper.CreateUsers(1, usersCollection)

				userId := user.InsertedIDs[0].(primitive.ObjectID).Hex()

				requestBody := map[string]string{
					"content": faker.Paragraph(),
				}

				endpoint := "/users/" + userId + "/posts"

				resp := helper.MakePostRequest(configs.GetString("app.fiber.address"), endpoint, requestBody)

				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusCreated))
			})
		})

		Context("when its given a quoted-post", func() {
			It("creates it without error", func() {
				user := helper.CreateUsers(1, usersCollection)

				userId := user.InsertedIDs[0].(primitive.ObjectID)

				post := helper.CreatePosts(1, userId, postsCollection)

				postId := post.InsertedIDs[0].(primitive.ObjectID).Hex()

				requestBody := map[string]string{
					"content":   faker.Paragraph(),
					"parent_id": postId,
				}

				endpoint := "/users/" + userId.Hex() + "/posts"

				resp := helper.MakePostRequest(configs.GetString("app.fiber.address"), endpoint, requestBody)

				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusCreated))
			})

		})

		Context("when its given a repost", func() {
			It("creates it without error", func() {
				user := helper.CreateUsers(1, usersCollection)

				userId := user.InsertedIDs[0].(primitive.ObjectID)

				post := helper.CreatePosts(1, userId, postsCollection)

				postId := post.InsertedIDs[0].(primitive.ObjectID).Hex()

				requestBody := map[string]string{
					"parent_id": postId,
				}

				endpoint := "/users/" + userId.Hex() + "/posts"

				resp := helper.MakePostRequest(configs.GetString("app.fiber.address"), endpoint, requestBody)

				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusCreated))
			})
		})

		Context("when its reach maximum number of posts per day", func() {
			It("returns error and not creates a new post", func() {
				user := helper.CreateUsers(1, usersCollection)

				userId := user.InsertedIDs[0].(primitive.ObjectID)

				helper.CreatePosts(configs.GetInt("app.posts.maximum-per-day"), userId, postsCollection)

				requestBody := map[string]string{
					"content": faker.Paragraph(),
				}

				endpoint := "/users/" + userId.Hex() + "/posts"

				resp := helper.MakePostRequest(configs.GetString("app.fiber.address"), endpoint, requestBody)

				Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusBadRequest))
			})
		})

		Describe("Getting last posts by user", func() {
			Context("when its given a user", func() {
				It("returns posts from this user", func() {
					expectedPostsReturned := 5

					user := helper.CreateUsers(1, usersCollection)

					userId := user.InsertedIDs[0].(primitive.ObjectID)

					helper.CreatePosts(10, userId, postsCollection)

					endpoint := "/users/" + userId.Hex() + "/posts"

					resp := helper.MakeGetRequest(configs.GetString("app.fiber.address"), endpoint, map[string]string{
						"page":  "1",
						"limit": strconv.Itoa(expectedPostsReturned),
					})

					var posts []entity.Post

					json.NewDecoder(resp.Body).Decode(&posts)

					Expect(len(posts)).To(Equal(expectedPostsReturned))
					Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusOK))
				})
			})
		})

		Describe("Getting user feed", func() {
			Context("when its given a user", func() {
				It("returns posts from its followers", func() {
					expectedPostsReturned := 10

					users := helper.CreateUsers(5, usersCollection)
					follower := users.InsertedIDs[0].(primitive.ObjectID)

					firstFollowed := users.InsertedIDs[1].(primitive.ObjectID)
					secondFollowed := users.InsertedIDs[2].(primitive.ObjectID)

					helper.CreateFollower(
						follower,
						firstFollowed,
						followersCollection,
					)

					helper.CreateFollower(
						follower,
						secondFollowed,
						followersCollection,
					)

					helper.CreatePosts(5, firstFollowed, postsCollection)
					helper.CreatePosts(5, secondFollowed, postsCollection)

					endpoint := "/users/" + follower.Hex() + "/feed"

					resp := helper.MakeGetRequest(configs.GetString("app.fiber.address"), endpoint, map[string]string{
						"page":  "1",
						"limit": strconv.Itoa(expectedPostsReturned),
					})

					var posts []entity.Post

					json.NewDecoder(resp.Body).Decode(&posts)

					Expect(len(posts)).To(Equal(expectedPostsReturned))
					Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusOK))

					for _, post := range posts {
						Expect([]primitive.ObjectID{firstFollowed, secondFollowed}).To(ContainElement(post.UserID))
					}
				})
			})

			Context("when its given a user that doesnt follow anyone", func() {
				It("returns no posts", func() {
					user := helper.CreateUsers(1, usersCollection)
					userId := user.InsertedIDs[0].(primitive.ObjectID)

					endpoint := "/users/" + userId.Hex() + "/feed"

					resp := helper.MakeGetRequest(configs.GetString("app.fiber.address"), endpoint, map[string]string{
						"page": "1",
					})

					var posts []entity.Post

					json.NewDecoder(resp.Body).Decode(&posts)

					Expect(resp.StatusCode).To(BeEquivalentTo(fiber.StatusOK))

					Expect(len(posts)).To(Equal(0))
				})
			})
		})

	})
})
