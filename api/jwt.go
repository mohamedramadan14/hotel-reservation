package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohamedramadan14/hotel-reservation/db"
	"net/http"
	"os"
	"time"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Println("NO TOKEN")
			//return fmt.Errorf("unauthorized: %s", "Token is not present")
			return ErrUnAuthorized()
		}
		claims, err := validateToken(token[len(token)-1])
		if err != nil {
			return err
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrUnAuthorized()
		}
		// set current authenticated user to the context
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method")
			return nil, ErrUnAuthorized()
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		log.Fatal("Failed to parse JWT token", err)
		return nil, NewError(http.StatusUnauthorized, "invalid credentials")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	expireStr, ok := claims["expire"].(string)
	if !ok || !token.Valid {
		return nil, NewError(http.StatusUnauthorized, "invalid credentials")
	}

	expireTime, err := time.Parse(time.RFC3339, expireStr)
	if err != nil {
		return nil, NewError(http.StatusUnauthorized, "token expired")
	}

	if time.Now().After(expireTime) {
		return nil, NewError(http.StatusUnauthorized, "token expired")
	}

	return claims, nil
}
