package main

import (
	"github.com/regiszanandrea/posty/internal"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		internal.ApplicationModule,
		internal.ApplicationInvokables,
	).Run()
}
