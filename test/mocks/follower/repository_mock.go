package follower_mock

import (
	follower_repository "github.com/regiszanandrea/posty/internal/user/repository/follower"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SuccessFollowerRepositoryMock struct {
	follower_repository.Repository
}

func (repo *SuccessFollowerRepositoryMock) GetFollowingUsers(followerID string) ([]string, error) {
	return []string{
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
	}, nil
}

func (repo *SuccessFollowerRepositoryMock) Follow(followerId, followingId string) (string, error) {
	return primitive.NewObjectID().Hex(), nil
}
func (repo *SuccessFollowerRepositoryMock) Unfollow(followerId, followingId string) error {
	return nil
}
