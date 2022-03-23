package fiber

import (
	"context"

	. "go.uber.org/fx"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

var (
	Module     = Provide(NewFiber)
	Invokables = Invoke(RegisterFiber)
)

func NewFiber(configs *viper.Viper) *fiber.App {
	return fiber.New(fiber.Config{DisableStartupMessage: configs.GetBool("app.fiber.disable-startup-message")})
}

func RegisterFiber(
	app *fiber.App,
	lifecycle Lifecycle,
	configs *viper.Viper,
) {
	lifecycle.Append(Hook{
		OnStart: func(context.Context) error {
			go app.Listen(configs.GetString("app.fiber.address")) //nolint
			return nil
		},
		OnStop: func(context.Context) error { return app.Shutdown() },
	})
}
