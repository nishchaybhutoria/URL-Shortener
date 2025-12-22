package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/nishchaybhutoria/URL-Shortener/db"
	"github.com/redis/go-redis/v9"
)

func ResolveURL(c fiber.Ctx) error {
	url := c.Params("url")

	r := db.CreateClient(0)
	defer r.Close()

	val, err := r.Get(db.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL does not exist",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot connect to DB",
		})
	}

	return c.Redirect().Status(fiber.StatusMovedPermanently).To(val)
}
