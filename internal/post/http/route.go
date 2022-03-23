package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/regiszanandrea/posty/internal/post/http/handler"
)

func RegisterPostRoutes(
	app *fiber.App,
	postCreatorHandler *handler.PostCreatorHandler,
	postListerHandler *handler.PostListerHandler,
	feedListerHandler *handler.FeedListerHandler,
) {

	group := app.Group("/users/:id")

	group.Get("/feed", feedListerHandler.ListFeed)

	groupPost := group.Group("/posts")

	groupPost.Post("/", postCreatorHandler.CreatePost)
	groupPost.Get("/", postListerHandler.ListLastPosts)
}
