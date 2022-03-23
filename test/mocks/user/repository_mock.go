package user_mock

import (
	"errors"
	"github.com/regiszanandrea/posty/internal/user/entity"
	"github.com/regiszanandrea/posty/internal/user/repository/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"
)

type SuccessUserRepositoryMock struct {
	user_repository.Repository
}

func (repo *SuccessUserRepositoryMock) Find(id string) (*entity.User, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)

	return &entity.User{
		ID:             objectId,
		Username:       "test",
		CreatedAt:      time.Now(),
		FollowersCount: uint(rand.Uint64()),
		FollowingCount: uint(rand.Uint64()),
		PostsCount:     uint(rand.Uint64()),
	}, nil
}

func (repo *SuccessUserRepositoryMock) Create(user *entity.User) (string, error) {
	return primitive.NewObjectID().Hex(), nil
}

func (repo *SuccessUserRepositoryMock) IncrementFollowers(id string) error {
	return nil
}
func (repo *SuccessUserRepositoryMock) IncrementFollowing(id string) error {
	return nil
}
func (repo *SuccessUserRepositoryMock) DecrementFollowers(id string) error {
	return nil
}
func (repo *SuccessUserRepositoryMock) DecrementFollowing(id string) error {
	return nil
}

func (repo *SuccessUserRepositoryMock) IncreasePostsCount(id string) error {
	return nil
}

type ErrorOnFindingUserRepositoryMock struct {
	user_repository.Repository
}

func (repo *ErrorOnFindingUserRepositoryMock) Find(id string) (*entity.User, error) {
	return nil, errors.New("error on finding")
}
