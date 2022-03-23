package handler

import (
	. "go.uber.org/fx"
)

var (
	Module = Provide(
		NewUserFinderHandler,
		NewUserCreatorHandler,
		NewFollowUserHandler,
		NewUnfollowUserHandler,
	)
)
