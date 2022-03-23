package handler

import (
	"github.com/regiszanandrea/posty/internal/user/entity"
	"github.com/regiszanandrea/posty/internal/user/service"

	"github.com/gofiber/fiber/v2"
)

type FollowUserHandler struct {
	service service.Service
}

func NewFollowUserHandler(s service.Service) *FollowUserHandler {
	return &FollowUserHandler{
		service: s,
	}
}

func (h *FollowUserHandler) FollowUser(ctx *fiber.Ctx) error {

	err := h.service.Follow(&entity.FollowRequest{
		FollowerID:  ctx.Params("followerId"),
		FollowingID: ctx.Params("userId"),
	})

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	return ctx.JSON(fiber.Map{"message": "user followed with success"})
}
