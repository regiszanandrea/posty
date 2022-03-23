package internal

import (
	"github.com/regiszanandrea/posty/configs/app"
	"github.com/regiszanandrea/posty/internal/fiber"
	"github.com/regiszanandrea/posty/internal/mongodb"
	"github.com/regiszanandrea/posty/internal/post"
	"github.com/regiszanandrea/posty/internal/user"

	. "go.uber.org/fx"
)

var (
	ApplicationModule = Options(
		app.Module,
		fiber.Module,
		Provide(mongodb.NewMongoDBClient),
		user.Module,
		post.Module,
	)

	ApplicationInvokables = Options(
		fiber.Invokables,
		Invoke(mongodb.RegisterMongoDB),
		user.Invokables,
		post.Invokables,
	)
)
