package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mohamedramadan14/hotel-reservation/db"
	"github.com/mohamedramadan14/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) tearDown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
func setup(_ *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, db.TestDBNAME),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)
	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "testFirstName",
		LastName:  "testLastName",
		Email:     "test@test.com",
		Password:  "test123456",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return
	}
	if user.ID == primitive.NilObjectID {
		t.Errorf("expected %s, got %s", "not nil", user.ID)
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Expected EncryptedPassword not to be set")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected %s, got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected %s, got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected %s, got %s", params.Email, user.Email)
	}
}
