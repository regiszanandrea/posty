package handler

import (
	"github.com/regiszanandrea/posty/internal/user/entity"
	"github.com/regiszanandrea/posty/internal/user/service"

	"github.com/gofiber/fiber/v2"
)

type UserCreatorHandler struct {
	service service.Service
}

func NewUserCreatorHandler(s service.Service) *UserCreatorHandler {
	return &UserCreatorHandler{
		service: s,
	}
}

func (h UserCreatorHandler) CreateUser(ctx *fiber.Ctx) error {
	user := new(entity.User)

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	id, errors := h.service.CreateUser(user)

	if errors != nil {
		var errorsStr []string

		for _, e := range errors {
			errorsStr = append(errorsStr, e.Error())
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(errorsStr)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}
