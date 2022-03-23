package entity

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Username       string             `bson:"username" validate:"required,max=14,alphanum"`
	CreatedAt      time.Time          `bson:"created_at"`
	FollowersCount uint               `bson:"followers_count"`
	FollowingCount uint               `bson:"following_count"`
	PostsCount     uint               `bson:"posts_count"`
}

type Follower struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	FollowerID  primitive.ObjectID `bson:"follower_id,omitempty"`
	FollowingID primitive.ObjectID `bson:"user_id,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"`
}

type FollowRequest struct {
	FollowerID  string `json:"follower_id" validate:"required"`
	FollowingID string `json:"user_id" validate:"required"`
}

type UnfollowRequest struct {
	FollowerID  string `json:"follower_id" validate:"required"`
	FollowingID string `json:"user_id" validate:"required"`
}

var validate = validator.New()

func Validate(user *User) []error {
	var errs []error
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, errors.New("field: "+err.StructField()+" "+err.Tag()))
		}
	}

	return errs
}
