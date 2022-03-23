package http

import (
	. "go.uber.org/fx"
)

var (
	Invokables = Invoke(
		RegisterPostRoutes,
	)
)
