package service

import (
	"errors"
	"github.com/regiszanandrea/posty/internal/post/entity"
	"github.com/regiszanandrea/posty/internal/post/repository"
	follower_repository "github.com/regiszanandrea/posty/internal/user/repository/follower"
	"github.com/spf13/viper"
	"time"
)

var ErrLimitPostsByDay = errors.New("limit of posts by day reached")

type Service interface {
	CreatePost(createPostRequest *entity.CreatePostRequest) (*string, []error)
	ListLastPostByUser(listPostRequest *entity.ListPostRequest) ([]*entity.Post, []error)
	ListFeed(listFeedRequest *entity.ListFeedRequest) ([]*entity.Post, []error)
}

type PostService struct {
	repository         post_repository.Repository
	followerRepository follower_repository.Repository
	configs            *viper.Viper
}

func NewPostService(
	repository post_repository.Repository,
	followerRepository follower_repository.Repository,
	configs *viper.Viper,
) *PostService {
	return &PostService{
		repository:         repository,
		followerRepository: followerRepository,
		configs:            configs,
	}
}

func (service *PostService) CreatePost(createPostRequest *entity.CreatePostRequest) (*string, []error) {
	errs := entity.Validate(createPostRequest)

	if errs != nil {
		return nil, errs
	}

	postsNumber, err := service.repository.GetNumberOfUsersPostsByDay(createPostRequest.UserID, time.Now())

	if err != nil {
		return nil, []error{err}
	}

	if postsNumber >= service.configs.GetInt("app.posts.maximum-per-day") {
		return nil, []error{ErrLimitPostsByDay}
	}

	id, err := service.repository.Create(createPostRequest.UserID, createPostRequest.ParentID, createPostRequest.Content)

	if err != nil {
		return nil, []error{err}
	}

	return &id, nil
}

func (service *PostService) ListLastPostByUser(listPostRequest *entity.ListPostRequest) ([]*entity.Post, []error) {
	errs := entity.ValidateStruct(listPostRequest)

	if errs != nil {
		return nil, errs
	}

	if listPostRequest.Limit == 0 {
		listPostRequest.Limit = service.configs.GetInt("app.posts.list-user-posts-limit")
	}

	posts, err := service.repository.GetLastByUser(listPostRequest.UserID, listPostRequest.Page, listPostRequest.Limit)

	if err != nil {
		return nil, []error{err}
	}

	return posts, nil
}

func (service *PostService) ListFeed(listFeedRequest *entity.ListFeedRequest) ([]*entity.Post, []error) {
	errs := entity.ValidateStruct(listFeedRequest)

	if errs != nil {
		return nil, errs
	}

	if listFeedRequest.Limit == 0 {
		listFeedRequest.Limit = service.configs.GetInt("app.posts.feed-posts-limit")
	}

	users, err := service.followerRepository.GetFollowingUsers(listFeedRequest.UserID)

	if len(users) == 0 {
		return nil, nil
	}

	posts, err := service.repository.GetLastByUsers(users, listFeedRequest.Page, listFeedRequest.Limit)

	if err != nil {
		return nil, []error{err}
	}

	return posts, nil
}
