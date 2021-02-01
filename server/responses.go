package server

import "github.com/gofiber/fiber/v2"

func HttpError(ctx *fiber.Ctx, status int, resp interface{}) error {
	ctx.Status(status)
	return ctx.JSON(resp)
}
