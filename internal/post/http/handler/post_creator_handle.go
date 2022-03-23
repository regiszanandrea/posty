package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/regiszanandrea/posty/internal/post/entity"
	"github.com/regiszanandrea/posty/internal/post/service"
	userService "github.com/regiszanandrea/posty/internal/user/service"
)

type PostCreatorHandler struct {
	service     service.Service
	userService userService.Service
}

func NewPostCreatorHandler(s service.Service, us userService.Service) *PostCreatorHandler {
	return &PostCreatorHandler{
		service:     s,
		userService: us,
	}
}

func (h *PostCreatorHandler) CreatePost(ctx *fiber.Ctx) error {
	post := new(entity.CreatePostRequest)

	if err := ctx.BodyParser(post); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	post.UserID = ctx.Params("id")

	id, errors := h.service.CreatePost(post)

	if errors != nil {
		var errorsStr []string

		for _, e := range errors {
			errorsStr = append(errorsStr, e.Error())
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(errorsStr)
	}

	err := h.userService.IncreaseNumberOfPosts(post.UserID)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}
