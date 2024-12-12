package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"go-service-template/internal/apperr"
	"go-service-template/internal/applog"
)

var v = validator.New()

func parseAndValidate(ctx *fiber.Ctx, req interface{}) error {
	err := ctx.BodyParser(req)
	if err != nil {
		applog.Logger.Error(ctx.UserContext(), err, "failed to parse the body", map[string]interface{}{
			"body": string(ctx.Body()),
		})

		return apperr.New(http.StatusInternalServerError, err, "failed to parse the body", apperr.ErrInternalError)
	}

	err = v.Struct(req)
	if err != nil {
		applog.Logger.Error(ctx.UserContext(), err, "failed to validate the request", map[string]interface{}{
			"request": req,
		})

		return err
	}

	return nil
}
