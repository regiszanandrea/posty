package post_repository

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/regiszanandrea/posty/configs/app"
	"github.com/regiszanandrea/posty/internal/mongodb"
	"github.com/regiszanandrea/posty/internal/post/entity"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestFollowerRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PostRepository Suite")
}

var (
	postRepository *PostRepository
	client         *mongo.Client
	configs        *viper.Viper
)

var _ = BeforeSuite(func() {
	configs = app.RegisterAppConfigs()

	client = mongodb.NewMongoDBClient(configs)

	err := client.Connect(context.Background())

	if err != nil {
		panic(err)
	}

	err = mongodb.CreateIndexes(client, configs, context.Background())

	if err != nil {
		panic(err)
	}

	postRepository = NewPostRepository(
		client, configs,
	)
})

var _ = AfterSuite(func() {
	postsCollection := client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.post-collection"),
	)

	_, err := postsCollection.DeleteMany(context.Background(), bson.M{})

	if err != nil {
		panic(err)
	}
})

var _ = Describe("PostRepository suite test", func() {
	Describe("Creating a post", func() {
		Context("when its given a valid post", func() {
			It("creates it without error", func() {
				id, err := postRepository.Create(primitive.NewObjectID().Hex(), "this is a post", "")

				Expect(err).To(BeNil())
				Expect(primitive.IsValidObjectID(id)).To(BeTrue())
			})
		})

		Context("when its given a quote-post", func() {
			It("creates it without error", func() {
				request := entity.CreatePostRequest{
					UserID:  primitive.NewObjectID().Hex(),
					Content: "this is a post",
				}

				id, _ := postRepository.Create(request.UserID, request.Content, "")

				parentId, _ := primitive.ObjectIDFromHex(id)

				quotePost := entity.CreatePostRequest{
					UserID:   primitive.NewObjectID().Hex(),
					Content:  "this is a quote-post",
					ParentID: parentId.Hex(),
				}

				id, err := postRepository.Create(quotePost.UserID, quotePost.Content, quotePost.ParentID)

				Expect(err).To(BeNil())
				Expect(primitive.IsValidObjectID(id)).To(BeTrue())
			})
		})
	})

	Describe("Getting number of posts by day", func() {
		Context("when its given two posts on a day", func() {
			It("returns two posts", func() {
				numberOfPostsExpected := 2

				user := primitive.NewObjectID()
				createPosts(user, numberOfPostsExpected)

				numberOfPosts, err := postRepository.GetNumberOfUsersPostsByDay(user.Hex(), time.Now())

				Expect(err).To(BeNil())
				Expect(numberOfPosts).To(Equal(numberOfPostsExpected))
			})
		})

		Context("when its given a user without posts on a day", func() {
			It("returns no posts", func() {
				_, _ = postRepository.Create(primitive.NewObjectID().Hex(), "this is a post", "")

				numberOfPosts, err := postRepository.GetNumberOfUsersPostsByDay(primitive.NewObjectID().Hex(), time.Now())

				Expect(err).To(BeNil())
				Expect(numberOfPosts).To(Equal(0))
			})
		})

		Context("when its given a user with posts but not at that day", func() {
			It("returns no posts", func() {
				_, _ = postRepository.Create(primitive.NewObjectID().Hex(), "this is a post", "")

				numberOfPosts, err := postRepository.GetNumberOfUsersPostsByDay(
					primitive.NewObjectID().Hex(),
					time.Now().AddDate(0, 0, 3),
				)

				Expect(err).To(BeNil())
				Expect(numberOfPosts).To(Equal(0))
			})
		})
	})

	Describe("Getting last posts by user", func() {
		Context("when its given five posts", func() {
			It("returns all posts", func() {
				numberOfPostsExpected := 5
				page := 1
				limit := 5

				user := primitive.NewObjectID()
				createPosts(user, numberOfPostsExpected)

				posts, err := postRepository.GetLastByUser(user.Hex(), page, limit)

				Expect(err).To(BeNil())
				Expect(len(posts)).To(Equal(numberOfPostsExpected))
			})
		})

		Context("when its given posts from other users", func() {
			It("returns only posts from the user specified", func() {
				numberOfPostsExpected := 3
				page := 1
				limit := 5

				user := primitive.NewObjectID()
				createPosts(user, numberOfPostsExpected)

				anotherUser := primitive.NewObjectID()
				createPosts(anotherUser, 2)

				anotherUser = primitive.NewObjectID()
				createPosts(anotherUser, 2)

				posts, err := postRepository.GetLastByUser(user.Hex(), page, limit)

				Expect(err).To(BeNil())
				Expect(len(posts)).To(Equal(numberOfPostsExpected))

				for _, post := range posts {
					Expect(post.UserID).To(Equal(user))
				}
			})
		})

		Context("when its given posts from other users", func() {
			It("returns only posts from the user specified", func() {
				numberOfPostsExpected := 3
				page := 1
				limit := 5

				user := primitive.NewObjectID()
				createPosts(user, numberOfPostsExpected)

				anotherUser := primitive.NewObjectID()
				createPosts(anotherUser, 2)

				anotherUser = primitive.NewObjectID()
				createPosts(anotherUser, 2)

				posts, err := postRepository.GetLastByUser(user.Hex(), page, limit)

				Expect(err).To(BeNil())
				Expect(len(posts)).To(Equal(numberOfPostsExpected))

				for _, post := range posts {
					Expect(post.UserID).To(Equal(user))
				}
			})
		})

		Context("when its given search by second page of getting last users posts", func() {
			It("returns only posts from the second page", func() {
				numberOfPostsExpected := 10
				page := 2
				limit := 5

				user := primitive.NewObjectID()
				createPosts(user, numberOfPostsExpected)

				posts, err := postRepository.GetLastByUser(user.Hex(), page, limit)

				Expect(err).To(BeNil())
				Expect(len(posts)).To(Equal(limit))
			})
		})

	})

	Describe("Getting last posts by users", func() {
		Context("when its given search by posts from users", func() {
			It("returns only posts from the users", func() {
				numberOfPosts := 2
				numberOfUsers := 10
				page := 2
				limit := 5

				var users []string
				for i := 0; i < numberOfUsers; i++ {
					user := primitive.NewObjectID()
					users = append(users, user.Hex())
					createPosts(user, numberOfPosts)
				}

				posts, err := postRepository.GetLastByUsers(users, page, limit)

				Expect(err).To(BeNil())
				Expect(len(posts)).To(Equal(limit))

				for _, post := range posts {
					Expect(contains(users, post.UserID.Hex())).To(BeTrue())
				}
			})
		})
	})
})

func createPosts(user primitive.ObjectID, numberOfPosts int) []string {
	var posts []string
	for i := 0; i < numberOfPosts; i++ {
		post, _ := postRepository.Create(user.Hex(), "", "this is a post")

		posts = append(posts, post)
	}
	return posts
}

func contains(items []string, needle string) bool {
	for _, item := range items {
		if item == needle {
			return true
		}
	}
	return false
}
