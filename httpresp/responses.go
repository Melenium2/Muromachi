package httpresp

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func Error(ctx *fiber.Ctx, status int, resp interface{}) error {
	ctx.Status(status)
	switch v := resp.(type) {
	case error, string:
		return ctx.JSON(map[string]interface{}{
			"error": fmt.Sprintf("%v", v),
			"status": status,
		})
	default:
		return ctx.JSON(resp)
	}
}
