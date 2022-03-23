package post

import (
	"github.com/regiszanandrea/posty/internal/post/http"
	"github.com/regiszanandrea/posty/internal/post/http/handler"
	"github.com/regiszanandrea/posty/internal/post/repository"
	"github.com/regiszanandrea/posty/internal/post/service"
	. "go.uber.org/fx"
)

var (
	Module = Options(
		Provide(
			Annotate(
				post_repository.NewPostRepository,
				As(new(post_repository.Repository)),
			),
			Annotate(
				service.NewPostService,
				As(new(service.Service)),
			),
		),
		handler.Module,
	)

	Invokables = Options(
		http.Invokables,
	)
)
