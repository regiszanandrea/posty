package handler

import (
	. "go.uber.org/fx"
)

var (
	Module = Provide(
		NewPostCreatorHandler,
		NewPostListerHandler,
		NewFeedListerHandler,
	)
)
