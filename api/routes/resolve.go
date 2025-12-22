package routes

import "github.com/gofiber/fiber/v3"

func ResolveURL(c fiber.Ctx) error {
	return c.SendString("resolving")
}
