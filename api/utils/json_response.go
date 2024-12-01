package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// JSONSuccess is a generic success output writer
func JSONSuccess(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(data)
}

// JSONFail is a generic fail output writer
// JSONFail can used for 4xx status code response
func JSONFail(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(data)
}

// JSONError is a generic error output writer
// JSONError can used for 5xx status code response
func JSONError(c *fiber.Ctx, statusCode int, err string) error {
	return c.Status(statusCode).JSON(errors.New(err))
}
