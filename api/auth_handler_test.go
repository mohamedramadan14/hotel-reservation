package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedramadan14/hotel-reservation/db/fixtures"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	testUser := fixtures.AddUser(tdb.Store, "mo", "gmail", false)
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthentication)

	params := AuthParams{
		Email:    "mo@gmail.com",
		Password: "mo_gmail",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Error(err)
	}

	if authResp.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}
	// Set encrypted Password to empty string because we don't want to return password in JSON Response
	testUser.EncryptedPassword = ""
	if !reflect.DeepEqual(testUser, authResp.User) {
		fmt.Println(testUser, authResp.User)
		t.Fatalf("expected user to match one with token")
	}
}

func TestAuthenticateWithWrongPasswordFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	fixtures.AddUser(tdb.Store, "mo", "gmail", false)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthentication)

	params := AuthParams{
		Email:    "mo@gmail.com",
		Password: "testPass_",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected http status of 401 but got %d", resp.StatusCode)
	}
	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected Type to be 'error' but got %s", genResp.Type)
	}

	if genResp.Msg != "invalid Credentials" {
		t.Fatalf("expected Msg to be 'invalid Credentials' but got %s", genResp.Msg)
	}

}
