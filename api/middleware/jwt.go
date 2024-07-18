package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	if err := validateToken(token[len(token)-1]); err != nil {
		return err
	}
	return c.Next()
}

func validateToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method")
			return nil, fmt.Errorf("unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		log.Fatal("Failed to parse JWT token", err)
		return fmt.Errorf("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid || time.Unix(int64(claims["expire"].(float64)), 0).Before(time.Now()) {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}
