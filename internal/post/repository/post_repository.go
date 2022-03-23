package post_repository

import (
	"context"
	"github.com/regiszanandrea/posty/internal/mongodb"
	"github.com/regiszanandrea/posty/internal/post/entity"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	Create(userId, content, parentId string) (string, error)
	GetLastByUser(userId string, page, limit int) ([]*entity.Post, error)
	GetLastByUsers(users []string, page, limit int) ([]*entity.Post, error)
	GetNumberOfUsersPostsByDay(id string, day time.Time) (int, error)
}

type PostRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(client *mongo.Client, configs *viper.Viper) *PostRepository {
	postsCollection := client.Database(
		configs.GetString("app.mongodb.database"),
	).Collection(
		configs.GetString("app.mongodb.post-collection"),
	)

	return &PostRepository{
		collection: postsCollection,
	}
}

func (repo *PostRepository) Create(userId, parentId, content string) (string, error) {
	objectId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return "", err
	}

	post := entity.Post{
		UserID:    objectId,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if parentId != "" {
		parentId, err := primitive.ObjectIDFromHex(userId)

		if err != nil {
			return "", err
		}

		post.ParentID = parentId
	}

	result, err := repo.collection.InsertOne(context.TODO(), post)

	if err != nil {
		if mongodb.IsDup(err) {
			return "", mongodb.ErrDuplicateKey
		}
		return "", err
	}

	objectId = result.InsertedID.(primitive.ObjectID)

	return objectId.Hex(), err
}

func (repo *PostRepository) GetLastByUser(userId string, page, limit int) ([]*entity.Post, error) {
	objectId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, err
	}

	var result []*entity.Post

	matchStage := bson.D{{"$match", bson.D{{"user_id", objectId}}}}

	sortStage := generateSortStage()
	paginationStage := generatePaginationStage(page, limit)
	limitStage := generateLimitStage(limit)
	lookUpStage := generateLookUpStage()
	quotedPostStage := generateQuotedPostStage()

	curr, err := repo.collection.Aggregate(context.TODO(), mongo.Pipeline{
		matchStage,
		sortStage,
		paginationStage,
		limitStage,
		lookUpStage,
		quotedPostStage,
	})

	if err != nil {
		return nil, err
	}

	for curr.Next(context.TODO()) {
		var post entity.Post
		if err := curr.Decode(&post); err != nil {
			return nil, err
		}

		result = append(result, &post)
	}

	return result, nil
}

func (repo *PostRepository) GetLastByUsers(users []string, page, limit int) ([]*entity.Post, error) {
	var result []*entity.Post

	var usersObjectId []primitive.ObjectID

	for _, user := range users {
		objId, err := primitive.ObjectIDFromHex(user)
		if err != nil {
			return nil, err
		}

		usersObjectId = append(usersObjectId, objId)
	}

	matchStage := bson.D{
		{
			"$match",
			bson.D{
				{
					"user_id",
					bson.D{
						{"$in", usersObjectId},
					},
				},
			},
		},
	}

	sortStage := generateSortStage()
	paginationStage := generatePaginationStage(page, limit)
	limitStage := generateLimitStage(limit)
	lookUpStage := generateLookUpStage()
	quotedPostStage := generateQuotedPostStage()

	curr, err := repo.collection.Aggregate(context.TODO(), mongo.Pipeline{
		matchStage,
		sortStage,
		paginationStage,
		limitStage,
		lookUpStage,
		quotedPostStage,
	})

	if err != nil {
		return nil, err
	}

	for curr.Next(context.TODO()) {
		var post entity.Post
		if err := curr.Decode(&post); err != nil {
			return nil, err
		}

		result = append(result, &post)
	}

	return result, nil
}

func (repo *PostRepository) GetNumberOfUsersPostsByDay(id string, day time.Time) (int, error) {
	day = time.Date(day.Year(), day.Month(), day.Day(), 6, 0, 0, day.Nanosecond(), day.Location())

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return 0, err
	}

	postNumber, err := repo.collection.CountDocuments(context.TODO(), bson.M{"user_id": objectId, "created_at": bson.M{
		"$gte": primitive.NewDateTimeFromTime(day),
	}})

	if err != nil {
		return 0, err
	}

	return int(postNumber), nil
}

func generateSortStage() bson.D {
	return bson.D{
		{"$sort", bson.D{{"created_at", -1}}},
	}
}

func generatePaginationStage(page int, limit int) bson.D {
	skip := int64(page*limit - limit)
	return bson.D{
		{"$skip", skip},
	}
}

func generateLimitStage(limit int) bson.D {
	return bson.D{
		{"$limit", limit},
	}
}

func generateLookUpStage() bson.D {
	return bson.D{
		{
			"$lookup",
			bson.D{
				{"from", "posts"},
				{"localField", "parent_id"},
				{"foreignField", "_id"},
				{"as", "quoted_post"},
			},
		},
	}
}

func generateQuotedPostStage() bson.D {
	return bson.D{
		{
			"$addFields",
			bson.D{
				{
					"quoted_post",
					bson.D{{"$first", "$quoted_post"}},
				},
			},
		},
	}
}
