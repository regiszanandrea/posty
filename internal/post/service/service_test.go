package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/regiszanandrea/posty/configs/app"
	"github.com/regiszanandrea/posty/internal/post/entity"
	follower_mock "github.com/regiszanandrea/posty/test/mocks/follower"
	"github.com/regiszanandrea/posty/test/mocks/post"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFollowerRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Post Service Suite")
}

var (
	configs *viper.Viper
	service *PostService
)

var _ = BeforeSuite(func() {
	configs = app.RegisterAppConfigs()
})

var _ = Describe("PostRepository suite test", func() {

	Describe("Creating a post", func() {
		BeforeEach(func() {
			service = NewPostService(
				&post_mock.SuccessPostRepositoryMock{},
				&follower_mock.SuccessFollowerRepositoryMock{},
				configs,
			)
		})

		Context("when its given a valid post", func() {
			It("creates it without error", func() {
				request := entity.CreatePostRequest{
					UserID:  primitive.NewObjectID().Hex(),
					Content: "this is a post",
				}

				id, err := service.CreatePost(&request)

				Expect(err).To(BeNil())
				Expect(primitive.IsValidObjectID(*id)).To(BeTrue())
			})
		})

		Context("when its given a quoted-post", func() {
			It("creates it without error", func() {
				request := entity.CreatePostRequest{
					UserID:   primitive.NewObjectID().Hex(),
					Content:  "this is a post",
					ParentID: primitive.NewObjectID().Hex(),
				}

				id, err := service.CreatePost(&request)

				Expect(err).To(BeNil())
				Expect(primitive.IsValidObjectID(*id)).To(BeTrue())
			})
		})

		Context("when its given a repost", func() {
			It("creates it without error", func() {
				request := entity.CreatePostRequest{
					UserID:   primitive.NewObjectID().Hex(),
					ParentID: primitive.NewObjectID().Hex(),
				}

				id, err := service.CreatePost(&request)

				Expect(err).To(BeNil())
				Expect(primitive.IsValidObjectID(*id)).To(BeTrue())
			})
		})

		Context("when its reach maximum number of posts per day", func() {
			It("returns error and not creates a new post", func() {
				service = NewPostService(
					&post_mock.MaximumPostsCreatedOnDayRepositoryMock{Configs: configs},
					&follower_mock.SuccessFollowerRepositoryMock{},
					configs,
				)

				request := entity.CreatePostRequest{
					UserID:  primitive.NewObjectID().Hex(),
					Content: "this is a post",
				}

				_, errors := service.CreatePost(&request)

				Expect(errors[0]).To(Equal(ErrLimitPostsByDay))
			})
		})
	})

	Describe("Getting last posts by user", func() {
		BeforeEach(func() {
			service = NewPostService(
				&post_mock.SuccessPostRepositoryMock{},
				&follower_mock.SuccessFollowerRepositoryMock{},
				configs,
			)
		})

		Context("when its given a user", func() {
			It("returns posts from this user", func() {
				userId := primitive.NewObjectID().Hex()
				request := entity.ListPostRequest{
					UserID: userId,
					Page:   1,
					Limit:  5,
				}

				posts, err := service.ListLastPostByUser(&request)

				Expect(err).To(BeNil())

				for _, post := range posts {
					Expect(post.UserID.Hex()).To(Equal(userId))
				}
			})
		})
	})

	Describe("Getting user feed", func() {
		BeforeEach(func() {
			service = NewPostService(
				&post_mock.SuccessPostRepositoryMock{},
				&follower_mock.SuccessFollowerRepositoryMock{},
				configs,
			)
		})

		Context("when its given a user", func() {
			It("returns posts from its followers", func() {
				userId := primitive.NewObjectID().Hex()
				request := entity.ListFeedRequest{
					UserID: userId,
					Page:   1,
					Limit:  5,
				}

				posts, err := service.ListFeed(&request)

				Expect(err).To(BeNil())

				for _, post := range posts {
					Expect(post.UserID.Hex()).NotTo(Equal(userId))
				}
			})
		})
	})
})
