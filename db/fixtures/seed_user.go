package fixtures

import (
	"context"
	"fmt"
	"log"

	"github.com/FerMusicComposer/hotel-reservation-backend/db"
	"github.com/FerMusicComposer/hotel-reservation-backend/models"
)

func seedUser(fname, lname, email, password, role string, userStore db.UserStore, ctx context.Context) {
	user, err := models.NewUserFromParams(models.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
		Role:      role,
	})

	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted user: ", insertedUser)

}
