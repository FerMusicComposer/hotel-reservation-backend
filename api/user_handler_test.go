package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
	"github.com/gofiber/fiber/v2"
)

const (
	testmongoURI = "mongodb://localhost:27017"
	tesdtDbName  = "test"
)

type testdb struct {
	db.UserStore
}

func setup(t *testing.T) *testdb {
	client, err := db.NewMongoConnection(testmongoURI, tesdtDbName)
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := models.CreateUserParams{
		FirstName: "Leonardo",
		LastName:  "da Vinci",
		Email:     "H3XK1@example.com",
		Password:  "password8978",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var user models.User
	json.NewDecoder(resp.Body).Decode(&user)

	if user.FirstName != params.FirstName {
		t.Errorf("expected first name %s, but got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("expected last name %s, but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected email %s, but got %s", params.Email, user.Email)
	}

	if len(user.ID) == 0 {
		t.Errorf("expected user ID to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected encrypted password not to be included on the JSON response")
	}

	tdb.teardown(t)
}
