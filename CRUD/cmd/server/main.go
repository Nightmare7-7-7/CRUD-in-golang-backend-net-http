package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"crud/internal/db"
	"crud/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {

	//Inject env
	err := godotenv.Load("../../.env")
	port := os.Getenv("PORT")

	if err != nil {
		log.Fatal("No .env file found", err)
	}

	//Connect DB
	database, err := db.Connect()

	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	println("Database connected successfully")
	defer database.Close()

	// prepare DB (migrations)
	err = db.RunMigrations(database)

	if err != nil {
		log.Fatal("Migration failed", err)
	}
	log.Println("✅ Database prepared successfully")

	//Setup and Listen GO server
	app := http.NewServeMux()

	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>welcome home</h1>")
	})

	//Routes
	app.HandleFunc("POST /create-user", handlers.CreateUser(database))
	app.HandleFunc("POST /login", handlers.LoginUser(database))
	app.HandleFunc("GET /user/{id}", handlers.GetUser(database))
	app.HandleFunc("GET /users/all", handlers.GetAllUsers(database))
	app.HandleFunc("PATCH /user/update-email", handlers.UpdateEmail(database))
	app.HandleFunc("DELETE /user/delete-user", handlers.DeleteUser(database))

	//Listen
	fmt.Println("Server running on http://localhost:" + port)
	http.ListenAndServe(":"+port, app)

}
