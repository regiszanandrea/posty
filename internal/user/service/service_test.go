package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/regiszanandrea/posty/internal/user/entity"
	"github.com/regiszanandrea/posty/test/mocks/follower"
	"github.com/regiszanandrea/posty/test/mocks/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFollowerRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Service Suite")
}

var (
	service *UserService
)

var _ = Describe("UserService suite test", func() {

	Describe("Getting a user", func() {
		Context("when its given a user", func() {
			It("returns data from this user", func() {
				service = NewUserService(
					&user_mock.SuccessUserRepositoryMock{},
					&follower_mock.SuccessFollowerRepositoryMock{},
				)

				userId := primitive.NewObjectID().Hex()
				user, err := service.GetUser(userId)

				Expect(err).To(BeNil())
				Expect(userId).To(Equal(user.ID.Hex()))
				Expect(user.FollowersCount).NotTo(Equal(0))
			})
		})

		Context("when its given a invalid user", func() {
			It("returns error", func() {
				service = NewUserService(
					&user_mock.ErrorOnFindingUserRepositoryMock{},
					&follower_mock.SuccessFollowerRepositoryMock{},
				)

				_, err := service.GetUser(primitive.NewObjectID().Hex())

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("Creating a user", func() {
		Context("when its given a valid user", func() {
			It("creates it without error", func() {
				service = NewUserService(
					&user_mock.SuccessUserRepositoryMock{},
					&follower_mock.SuccessFollowerRepositoryMock{},
				)

				id, err := service.CreateUser(&entity.User{Username: "testd"})

				Expect(err).To(BeNil())
				Expect(primitive.IsValidObjectID(*id)).To(BeTrue())
			})
		})
	})

	Describe("Following a user", func() {
		BeforeEach(func() {
			service = NewUserService(
				&user_mock.SuccessUserRepositoryMock{},
				&follower_mock.SuccessFollowerRepositoryMock{},
			)
		})

		Context("when its given a valid request", func() {
			It("follows the user", func() {
				request := entity.FollowRequest{
					FollowingID: primitive.NewObjectID().Hex(),
					FollowerID:  primitive.NewObjectID().Hex(),
				}

				err := service.Follow(&request)

				Expect(err).To(BeNil())
			})
		})
		Context("when its given same user as follower and following", func() {
			It("returns a error", func() {
				userId := primitive.NewObjectID().Hex()
				request := entity.FollowRequest{
					FollowingID: userId,
					FollowerID:  userId,
				}

				err := service.Follow(&request)

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("Unfollowing a user", func() {
		BeforeEach(func() {
			service = NewUserService(
				&user_mock.SuccessUserRepositoryMock{},
				&follower_mock.SuccessFollowerRepositoryMock{},
			)
		})

		Context("when its given a valid request", func() {
			It("unfollows the user", func() {
				request := entity.UnfollowRequest{
					FollowingID: primitive.NewObjectID().Hex(),
					FollowerID:  primitive.NewObjectID().Hex(),
				}

				err := service.Unfollow(&request)

				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Incrementing user number of posts", func() {
		BeforeEach(func() {
			service = NewUserService(
				&user_mock.SuccessUserRepositoryMock{},
				&follower_mock.SuccessFollowerRepositoryMock{},
			)
		})

		Context("when increment the user's posts number", func() {
			It("increments without error", func() {
				err := service.IncreaseNumberOfPosts(primitive.NewObjectID().Hex())

				Expect(err).To(BeNil())
			})
		})
	})
})
