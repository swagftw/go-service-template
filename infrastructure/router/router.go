package router

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"go-service-template/internal/apperr"
	"go-service-template/internal/applog"
)

func NewRouter() *fiber.App {
	app := fiber.New()

	app.Use(recover.New(), helmet.New(), requestid.New())
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))
	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err == nil {
			return nil
		}

		defaultErr := &apperr.AppError{
			Code:   http.StatusInternalServerError,
			Msg:    "internal server error",
			Status: apperr.ErrInternalError,
		}

		if err != nil {
			var appErr *apperr.AppError
			if errors.As(err, &appErr) {
				return c.Status(appErr.Code).JSON(appErr)
			}

			var fiberErr *fiber.Error

			if errors.As(err, &fiberErr) {
				defaultErr.Msg = fiberErr.Message
				defaultErr.Code = fiberErr.Code
				defaultErr.Status = apperr.ErrUnknown

				return c.Status(fiberErr.Code).JSON(fiberErr)
			}

			var validationErr validator.ValidationErrors

			if errors.As(err, &validationErr) {
				defaultErr.Msg = validationErr.Error()
				defaultErr.Code = http.StatusBadRequest
				defaultErr.Status = apperr.ErrBadRequest

				return c.Status(http.StatusBadRequest).JSON(defaultErr)
			}
		}

		return c.Status(defaultErr.Code).JSON(defaultErr)
	})

	return app
}

func Start(port int, app *fiber.App) error {
	hostPort := net.JoinHostPort("", strconv.Itoa(port))

	err := app.Listen(hostPort)
	if err != nil {
		applog.Logger.Error(context.Background(), err, "failed to start server", map[string]interface{}{
			"port": port,
		})

		return apperr.New(http.StatusInternalServerError, err, "failed to start server", apperr.ErrInternalError)
	}

	return nil
}
