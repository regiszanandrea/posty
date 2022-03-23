package follower_repository

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/regiszanandrea/posty/configs/app"
	"github.com/regiszanandrea/posty/internal/mongodb"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestFollowerRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FollowerRepository Suite")
}

var (
	followerRepository *FollowerRepository
	client             *mongo.Client
	configs            *viper.Viper
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

	followerRepository = NewFollowerRepository(
		client, configs,
	)
})

var _ = AfterSuite(func() {
	followersCollection := client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.follower-collection"),
	)

	_, err := followersCollection.DeleteMany(context.Background(), bson.M{})

	if err != nil {
		panic(err)
	}
})

var _ = Describe("FollowerRepository suite test", func() {
	Describe("Following a User", func() {
		Context("when its given a follower and the followed user", func() {
			It("persist the following between the two users", func() {
				_, err := followerRepository.Follow(primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex())

				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Unfollowing a User", func() {
		Context("when its a follower and the followed user", func() {
			It("persist the unfollowing between the two users", func() {
				followerId := primitive.NewObjectID().Hex()
				userId := primitive.NewObjectID().Hex()

				_, err := followerRepository.Follow(followerId, userId)

				err = followerRepository.Unfollow(followerId, userId)

				Expect(err).To(BeNil())
			})
		})

		Context("when there is no following", func() {
			It("not returns a error", func() {
				err := followerRepository.Unfollow(primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex())

				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Getting user's followers", func() {
		Context("when there is followers", func() {
			It("returns all of them", func() {
				followerId := primitive.NewObjectID().Hex()

				numberOfUsersFollowed := 10
				for i := 0; i < numberOfUsersFollowed; i++ {
					followerRepository.Follow(followerId, primitive.NewObjectID().Hex())
				}

				users, err := followerRepository.GetFollowingUsers(followerId)

				Expect(err).To(BeNil())
				Expect(len(users)).To(Equal(numberOfUsersFollowed))
			})
		})

		Context("when there is no followers", func() {
			It("returns nothing", func() {
				followerId := primitive.NewObjectID().Hex()

				users, err := followerRepository.GetFollowingUsers(followerId)

				Expect(err).To(BeNil())
				Expect(users).To(BeEmpty())
			})
		})
	})
})
