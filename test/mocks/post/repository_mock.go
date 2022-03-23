package post_mock

import (
	"github.com/regiszanandrea/posty/internal/post/entity"
	"github.com/regiszanandrea/posty/internal/post/repository"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SuccessPostRepositoryMock struct {
	post_repository.Repository
}

func (repo *SuccessPostRepositoryMock) GetNumberOfUsersPostsByDay(id string, day time.Time) (int, error) {
	return 0, nil
}

func (repo *SuccessPostRepositoryMock) Create(userId, content, parentId string) (string, error) {
	return primitive.NewObjectID().Hex(), nil
}

func (repo *SuccessPostRepositoryMock) GetLastByUser(userId string, page, limit int) ([]*entity.Post, error) {
	objectId, _ := primitive.ObjectIDFromHex(userId)

	return []*entity.Post{
		{
			UserID:  objectId,
			Content: "this is a post",
		},
		{
			UserID:  objectId,
			Content: "this is a second post",
		},
	}, nil
}

func (repo *SuccessPostRepositoryMock) GetLastByUsers(users []string, page, limit int) ([]*entity.Post, error) {
	var posts []*entity.Post
	for _, user := range users {
		objectId, _ := primitive.ObjectIDFromHex(user)

		posts = append(posts, &entity.Post{
			UserID:  objectId,
			Content: "this is a post",
		})
	}

	return posts, nil
}

type MaximumPostsCreatedOnDayRepositoryMock struct {
	post_repository.Repository
	Configs *viper.Viper
}

func (repo *MaximumPostsCreatedOnDayRepositoryMock) GetNumberOfUsersPostsByDay(id string, day time.Time) (int, error) {
	return repo.Configs.GetInt("app.posts.maximum-per-day"), nil
}
