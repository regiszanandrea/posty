package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/regiszanandrea/posty/internal/post/entity"
	"github.com/regiszanandrea/posty/internal/post/service"
)

type PostListerHandler struct {
	service service.Service
}

func NewPostListerHandler(s service.Service) *PostListerHandler {
	return &PostListerHandler{
		service: s,
	}
}

func (h *PostListerHandler) ListLastPosts(ctx *fiber.Ctx) error {
	list := new(entity.ListPostRequest)

	if err := ctx.QueryParser(list); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	list.UserID = ctx.Params("id")

	posts, errors := h.service.ListLastPostByUser(list)

	if errors != nil {
		var errorsStr []string

		for _, e := range errors {
			errorsStr = append(errorsStr, e.Error())
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(errorsStr)
	}

	return ctx.JSON(posts)
}
