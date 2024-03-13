package main

import (
	"go-oauth2-impl/libs"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Golang OAuth2 Server")
	})

	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	app.Get("/auth/google-login", func(c *fiber.Ctx) error {
		url := conf.AuthCodeURL("not-implemented-yet")

		return c.Redirect(url)
	})

	app.Get("/auth/google-callback", func(c *fiber.Ctx) error {
		code := c.Query("code")

		token, err := conf.Exchange(c.Context(), code)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		profile, err := libs.ConvertToken(token.AccessToken)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.JSON(profile)
	})

	app.Listen(":8000")
}
