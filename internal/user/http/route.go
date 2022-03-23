package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/regiszanandrea/posty/internal/user/http/handler"
)

func RegisterUserRoutes(
	app *fiber.App,
	userFinderHandler *handler.UserFinderHandler,
	userCreatorHandler *handler.UserCreatorHandler,
	followUserHandler *handler.FollowUserHandler,
	unfollowUserHandler *handler.UnfollowUserHandler,
) {
	group := app.Group("/users")

	group.Get("/:id", userFinderHandler.FindUser)
	group.Post("/", userCreatorHandler.CreateUser)
	group.Post("/:followerId/follow/:userId", followUserHandler.FollowUser)
	group.Post("/:followerId/unfollow/:userId", unfollowUserHandler.UnfollowUser)
}
