package user

import (
	"github.com/regiszanandrea/posty/internal/user/http"
	"github.com/regiszanandrea/posty/internal/user/http/handler"
	"github.com/regiszanandrea/posty/internal/user/repository/follower"
	"github.com/regiszanandrea/posty/internal/user/repository/user"
	"github.com/regiszanandrea/posty/internal/user/service"

	. "go.uber.org/fx"
)

var (
	Module = Options(
		Provide(
			Annotate(
				user_repository.NewUserRepository,
				As(new(user_repository.Repository)),
			),
			Annotate(
				follower_repository.NewFollowerRepository,
				As(new(follower_repository.Repository)),
			),
			Annotate(
				service.NewUserService,
				As(new(service.Service)),
			),
		),
		handler.Module,
	)

	Invokables = Options(
		http.Invokables,
	)
)
