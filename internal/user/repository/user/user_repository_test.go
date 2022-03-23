package user_repository

import (
	"context"
	"github.com/regiszanandrea/posty/internal/mongodb"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/regiszanandrea/posty/configs/app"
	"github.com/regiszanandrea/posty/internal/user/entity"
)

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserRepository Suite")
}

var (
	userRepository *UserRepository
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

	userRepository = NewUserRepository(
		client, configs,
	)
})

var _ = AfterSuite(func() {
	usersCollection := client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.user-collection"),
	)

	_, err := usersCollection.DeleteMany(context.Background(), bson.M{})

	if err != nil {
		panic(err)
	}
})

var _ = Describe("UserRepository suite test", func() {
	Describe("Creating a User", func() {
		Context("when its given a valid user", func() {
			It("creates it without error", func() {
				user := entity.User{
					Username:       "george",
					CreatedAt:      time.Now(),
					FollowersCount: 0,
					FollowingCount: 0,
				}

				id, err := userRepository.Create(&user)

				Expect(err).To(BeNil())
				Expect(primitive.IsValidObjectID(id)).To(BeTrue())
			})
		})

		Context("when already is a user with same username", func() {
			It("returns error", func() {
				user := entity.User{
					Username:       "michael",
					CreatedAt:      time.Now(),
					FollowersCount: 0,
					FollowingCount: 0,
				}

				_, err := userRepository.Create(&user)

				_, err = userRepository.Create(&user)

				Expect(err).To(Equal(mongodb.ErrDuplicateKey))
			})
		})
	})

	Describe("Finding a User", func() {
		Context("when there is a user", func() {
			It("returns the correct user", func() {
				user := entity.User{
					Username:       "sd",
					CreatedAt:      time.Now(),
					FollowersCount: 0,
					FollowingCount: 0,
				}

				id, _ := userRepository.Create(&user)

				userFound, err := userRepository.Find(id)
				Expect(err).To(BeNil())
				Expect(userFound.Username).To(Equal(user.Username))
			})
		})
	})

	Describe("Increment User's followers", func() {
		Context("when increment the user's followers", func() {
			It("increments only by one", func() {
				var expectedFollowers uint = 6

				user := entity.User{
					Username:       "testIncrementFollowers",
					CreatedAt:      time.Now(),
					FollowersCount: 5,
					FollowingCount: 0,
				}

				id, _ := userRepository.Create(&user)

				_ = userRepository.IncrementFollowers(id)

				userFound, err := userRepository.Find(id)

				Expect(err).To(BeNil())
				Expect(userFound.FollowersCount).To(Equal(expectedFollowers))
			})
		})
	})

	Describe("Increment User's following number", func() {
		Context("when increment the user's following number", func() {
			It("increments only by one", func() {
				var expectedFollowing uint = 1

				user := entity.User{
					Username:       "testIncrementFollowing",
					CreatedAt:      time.Now(),
					FollowersCount: 0,
					FollowingCount: 0,
				}

				id, _ := userRepository.Create(&user)

				_ = userRepository.IncrementFollowing(id)

				userFound, err := userRepository.Find(id)

				Expect(err).To(BeNil())
				Expect(userFound.FollowingCount).To(Equal(expectedFollowing))
			})
		})
	})

	Describe("Decrement User's followers", func() {
		Context("when decrement the user's followers", func() {
			It("decrements only by one", func() {
				var expectedFollowers uint = 4

				user := entity.User{
					Username:       "testDecrementFollowers",
					CreatedAt:      time.Now(),
					FollowersCount: 5,
					FollowingCount: 0,
				}

				id, _ := userRepository.Create(&user)

				_ = userRepository.DecrementFollowers(id)

				userFound, err := userRepository.Find(id)

				Expect(err).To(BeNil())
				Expect(userFound.FollowersCount).To(Equal(expectedFollowers))
			})
		})
	})

	Describe("Decrement User's following number", func() {
		Context("when decrement the user's following number", func() {
			It("decrements only by one", func() {
				var expectedFollowing uint = 0

				user := entity.User{
					Username:       "testDecrementFollowing",
					CreatedAt:      time.Now(),
					FollowersCount: 0,
					FollowingCount: 1,
				}

				id, _ := userRepository.Create(&user)

				_ = userRepository.DecrementFollowing(id)

				userFound, err := userRepository.Find(id)

				Expect(err).To(BeNil())
				Expect(userFound.FollowingCount).To(Equal(expectedFollowing))
			})
		})
	})

	Describe("Increment User's posts number", func() {
		Context("when increment the user's posts number", func() {
			It("increment only by one", func() {
				var expectedPostsCount uint = 1

				user := entity.User{
					Username:       "testIncrementPosts",
					CreatedAt:      time.Now(),
					FollowersCount: 0,
					FollowingCount: 0,
					PostsCount:     0,
				}

				id, _ := userRepository.Create(&user)

				_ = userRepository.IncreasePostsCount(id)

				userFound, err := userRepository.Find(id)

				Expect(err).To(BeNil())
				Expect(userFound.PostsCount).To(Equal(expectedPostsCount))
			})
		})
	})
})
