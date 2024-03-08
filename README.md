# Hotel reservation backend

1st project from Anthony GG's Full Time Go Dev course

I have implemented the mongo client in a slightly different way, that is cleaner in my opinion.
Also, the room model, booking circuit and authorization are different, as I have included some minor business logic which was out of the scope of the course, and for authorization, I did a simple RBAC middleware, that looks similar to the one of the course, but checks the user role from the JWT token claims, instead of the role from the database. Also I did a custom error implementation which vastly defers from the one in the course

The Postman endpoints collection is included for anyone that wants to test the API

In order to run the API locally, you need to do the following:

- Make sure Go is installed in your system https://go.dev/doc/tutorial/getting-started#install
- Install MongoDB Compass and create a connection https://www.mongodb.com/products/tools/compass
- To install all necessary packages cd into the project directory and run `go download`. After it you can run a `go mod tidy` to fix any possible errors
- Create a .env file in the root directory and copy the project environment variables
- After ensuring Mongo Compass is up and running, cd into the /scripts directory and run `go run db_seed.go`. This will seed the database with test data
- Check the collections are created, then cd into the root directory and run `air` to start the server
- If you want to test the API, import the collection into Postman. You will have to authenticate a user first in order to reach the endpoints, and
  some endpoints are accessible only to admin, for which you will have to authenticate the user with the admin role

# Project environment variables

```
HTTP_LISTEN_ADDRESS=:5000
JWT_SECRET=somethingsupersecretthatNOBODYKNOWS
```

## Project outline

- users -> book room from an hotel
- admins -> going to check reservation/bookings
- Authentication and authorization -> JWT tokens
- Hotels -> CRUD API -> JSON
- Rooms -> CRUD API -> JSON
- Scripts -> database management -> seeding, migration

## Resources

### Mongodb driver

Documentation

```
https://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongodb client

```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber

Documentation

```
https://gofiber.io
```

Installing gofiber

```
go get github.com/gofiber/fiber/v2
```

## Docker

### Installing mongodb as a Docker container

```
docker run --name mongodb -d mongo:latest -p 27017:27017
```
