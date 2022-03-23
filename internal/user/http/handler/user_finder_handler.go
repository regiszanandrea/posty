package handler

import (
	"github.com/regiszanandrea/posty/internal/user/service"

	"github.com/gofiber/fiber/v2"
)

type UserFinderHandler struct {
	service service.Service
}

func NewUserFinderHandler(s service.Service) *UserFinderHandler {
	return &UserFinderHandler{
		service: s,
	}
}

func (h UserFinderHandler) FindUser(ctx *fiber.Ctx) error {

	user, err := h.service.GetUser(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if user == nil {
		return ctx.Status(fiber.StatusNotFound).JSON("no user found")
	}

	return ctx.JSON(user)
}
