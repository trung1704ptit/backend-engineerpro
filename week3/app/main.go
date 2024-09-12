package main

import (
	"fmt"
	"log"
	"net/http"

	"simple/controllers"
	"simple/models"
)

func main() {
	// Initialize the SQLite database
	models.InitDB("user.db")

	// Route for user sign-up
	http.HandleFunc("/signup", controllers.SignUpHandler)

	// Route for user sign-in
	http.HandleFunc("/signin", controllers.SignInHandler)

	// Route for updating user profile
	http.HandleFunc("/update", controllers.UpdateProfileHandler)

	// Route for changing profile image
	http.HandleFunc("/upload-profile-image", controllers.ChangeProfileImageHandler)

	// Start the server
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
