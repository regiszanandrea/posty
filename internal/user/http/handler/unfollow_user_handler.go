package handler

import (
	"github.com/regiszanandrea/posty/internal/user/entity"
	"github.com/regiszanandrea/posty/internal/user/service"

	"github.com/gofiber/fiber/v2"
)

type UnfollowUserHandler struct {
	service service.Service
}

func NewUnfollowUserHandler(s service.Service) *UnfollowUserHandler {
	return &UnfollowUserHandler{
		service: s,
	}
}

func (h *UnfollowUserHandler) UnfollowUser(ctx *fiber.Ctx) error {

	err := h.service.Unfollow(&entity.UnfollowRequest{
		FollowerID:  ctx.Params("followerId"),
		FollowingID: ctx.Params("userId"),
	})

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	return ctx.JSON(fiber.Map{"message": "user unfollowed with success"})
}
