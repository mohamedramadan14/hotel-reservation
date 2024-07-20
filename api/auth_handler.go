package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohamedramadan14/hotel-reservation/db"
	"github.com/mohamedramadan14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	userStore db.UserStore
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"message"`
}

func invalidCredentialsResponseGenerator(c *fiber.Ctx) error {
	return c.Status(http.StatusUnauthorized).JSON(
		genericResp{
			Type: "error",
			Msg:  "invalid Credentials",
		})
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *types.User `json:"user"`
}

func (h *AuthHandler) HandleAuthentication(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentialsResponseGenerator(c)
		}
		return err
	}

	_, err = types.IsValidPassword(user.EncryptedPassword, authParams.Password)
	if err != nil {
		return invalidCredentialsResponseGenerator(c)
	}

	token := GenerateTokenFromUser(user)
	fmt.Println("authenticated ->", user.FirstName)
	resp := AuthResponse{
		Token: token,
		User:  user,
	}
	return c.JSON(resp)
}

func GenerateTokenFromUser(user *types.User) string {
	now := time.Now()
	validTill := now.Add(time.Hour * 72)
	claims := jwt.MapClaims{
		"id":     user.ID,
		"email":  user.Email,
		"expire": validTill,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}
