package service

import (
	"errors"
	"github.com/regiszanandrea/posty/internal/user/entity"
	"github.com/regiszanandrea/posty/internal/user/repository/follower"
	"github.com/regiszanandrea/posty/internal/user/repository/user"
	"time"
)

type Service interface {
	GetUser(id string) (*entity.User, error)
	CreateUser(user *entity.User) (*string, []error)
	Follow(followRequest *entity.FollowRequest) error
	Unfollow(unfollowRequest *entity.UnfollowRequest) error
	IncreaseNumberOfPosts(id string) error
}

type UserService struct {
	userRepository     user_repository.Repository
	followerRepository follower_repository.Repository
}

func NewUserService(userRepo user_repository.Repository, followerRepo follower_repository.Repository) *UserService {
	return &UserService{
		userRepository:     userRepo,
		followerRepository: followerRepo,
	}
}

func (service *UserService) GetUser(id string) (*entity.User, error) {
	user, err := service.userRepository.Find(id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) CreateUser(user *entity.User) (*string, []error) {
	errs := entity.Validate(user)

	if errs != nil {
		return nil, errs
	}

	user.CreatedAt = time.Now()

	id, err := service.userRepository.Create(user)

	if err != nil {
		return nil, []error{err}
	}
	return &id, nil
}

func (service *UserService) Follow(followRequest *entity.FollowRequest) error {
	if followRequest.FollowerID == followRequest.FollowingID {
		return errors.New("A user cannot follow itself")
	}

	_, err := service.followerRepository.Follow(followRequest.FollowerID, followRequest.FollowingID)

	if err != nil {
		return err
	}

	err = service.userRepository.IncrementFollowers(followRequest.FollowingID)

	if err != nil {
		return err
	}

	err = service.userRepository.IncrementFollowing(followRequest.FollowerID)

	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) Unfollow(unfollowRequest *entity.UnfollowRequest) error {
	err := service.followerRepository.Unfollow(unfollowRequest.FollowerID, unfollowRequest.FollowingID)

	if err != nil {
		return err
	}

	err = service.userRepository.DecrementFollowers(unfollowRequest.FollowingID)

	if err != nil {
		return err
	}

	err = service.userRepository.DecrementFollowing(unfollowRequest.FollowerID)

	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) IncreaseNumberOfPosts(id string) error {
	return service.userRepository.IncreasePostsCount(id)
}
