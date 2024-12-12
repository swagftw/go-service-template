package http

import (
	"github.com/gofiber/fiber/v2"

	"go-service-template/domain"
)

type healthHandler struct {
	healthService domain.IHealthService
}

func InitHealthRouter(rtr fiber.Router, healthService domain.IHealthService) {
	handler := healthHandler{
		healthService: healthService,
	}

	rtr.Post("/ping", handler.ping)
}

func (h healthHandler) ping(ctx *fiber.Ctx) error {
	req := new(PingReq)

	if err := parseAndValidate(ctx, req); err != nil {
		return err
	}

	res, err := h.healthService.Ping(ctx.UserContext(), &domain.PingReq{
		Name: req.Name,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(&PingRes{GreetMessage: res.GreetMessage})
}
