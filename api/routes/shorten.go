package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/nishchaybhutoria/URL-Shortener/db"
	"github.com/nishchaybhutoria/URL-Shortener/helpers"
	"github.com/redis/go-redis/v9"
)

type request struct {
	URL         string        `json:"url"`    // original URL
	CustomShort string        `json:"short"`  // custom short URL (optional)
	Expiry      time.Duration `json:"expiry"` // time to expire: int (hours)
}

type response struct {
	URL             string        `json:"url"`              // original URL
	CustomShort     string        `json:"short"`            // final short URL
	Expiry          time.Duration `json:"expiry"`           // time to expire: int (hours)
	XRateRemaining  string        `json:"rate_limit"`       // number of API calls remaining: int
	XRateLimitReset time.Duration `json:"rate_limit_reset"` // time to rate limit reset: int
}

func ShortenURL(c fiber.Ctx) error {
	// enforce rate limits
	// connect to rate limit db (DB 1)
	r2 := db.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(db.Ctx, c.IP()).Result()

	if err == redis.Nil {
		_ = r2.Set(db.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err() // new user: give API_QUOTA requests for 30 minutes
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(db.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	body := new(request)

	// ensure request is serializable
	if err := c.Bind().Body(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// check if URL is valid
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	// prevent redirect chains: ban shortening our domain
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Domain not allowed",
		})
	}

	body.URL = helpers.EnforceHTTP(body.URL)

	// validation done
	id := uuid.New().String()[:6]
	if body.CustomShort != "" {
		id = body.CustomShort
	}

	r1 := db.CreateClient(0)
	defer r1.Close()

	val, _ = r1.Get(db.Ctx, id).Result()

	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL short already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24 // set default expiry time
	}

	err = r1.Set(db.Ctx, id, body.URL, body.Expiry*60*60*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  "",
		XRateLimitReset: 30,
	}

	r2.Decr(db.Ctx, c.IP()) // decrement rate limit

	val, _ = r2.Get(db.Ctx, c.IP()).Result()
	resp.XRateRemaining = val

	ttl, _ := r2.TTL(db.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	return c.Status(fiber.StatusOK).JSON(resp)
}
