package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/regiszanandrea/posty/internal/post/entity"
	"github.com/regiszanandrea/posty/internal/post/service"
)

type FeedListerHandler struct {
	service service.Service
}

func NewFeedListerHandler(s service.Service) *FeedListerHandler {
	return &FeedListerHandler{
		service: s,
	}
}

func (h *FeedListerHandler) ListFeed(ctx *fiber.Ctx) error {
	list := new(entity.ListFeedRequest)

	if err := ctx.QueryParser(list); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	list.UserID = ctx.Params("id")

	posts, errors := h.service.ListFeed(list)

	if errors != nil {
		var errorsStr []string

		for _, e := range errors {
			errorsStr = append(errorsStr, e.Error())
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(errorsStr)
	}

	return ctx.JSON(posts)
}
