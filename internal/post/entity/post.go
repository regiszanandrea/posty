package entity

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	ParentID   primitive.ObjectID `bson:"parent_id,omitempty"`
	QuotedPost *QuotedPost        `bson:"quoted_post,omitempty"`
	Content    string             `bson:"content,omitempty"`
	CreatedAt  time.Time          `bson:"created_at"`
}

type QuotedPost struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Content   string             `bson:"content,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
}

type CreatePostRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	ParentID string `json:"parent_id"`
	Content  string `json:"content" validate:"max=777"`
}

type ListPostRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Page   int    `json:"page" validate:"required"`
	Limit  int    `json:"limit"`
}

type ListFeedRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Page   int    `json:"page" validate:"required"`
	Limit  int    `json:"limit"`
}

var validate = validator.New()

func Validate(createPostRequest *CreatePostRequest) []error {
	var errs []error

	if createPostRequest.Content == "" && createPostRequest.ParentID == "" {
		errs = append(errs, errors.New("field content must be present when field parentID is not"))
	}

	err := validate.Struct(createPostRequest)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, errors.New("field: "+err.StructField()+" "+err.Tag()))
		}
	}

	return errs
}

func ValidateStruct(st interface{}) []error {
	var errs []error

	err := validate.Struct(st)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, errors.New("field: "+err.StructField()+" "+err.Tag()))
		}
	}

	return errs
}
