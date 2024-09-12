package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"simple/models"
)

// SignUpHandler handles user registration
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Store user in the database
	if err := models.CreateUser(user); err != nil {
		http.Error(w, "User already exists or cannot create user", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s signed up successfully\n", user.Username)
}

// SignInHandler handles user sign-in
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Fetch the user from the database
	user, err := models.GetUserByUsername(credentials.Username)
	if err != nil || user.Password != credentials.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "User %s signed in successfully\n", credentials.Username)
}

// UpdateProfileHandler allows the user to update their profile
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var updateUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if user exists
	user, err := models.GetUserByUsername(updateUser.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update user profile
	user.Email = updateUser.Email
	user.Age = updateUser.Age
	if err := models.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s profile updated successfully\n", updateUser.Username)
}

// ChangeProfileImageHandler allows the user to upload a profile image
func ChangeProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")

	// Check if user exists
	user, err := models.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Parse the multipart form for file upload
	err = r.ParseMultipartForm(10 << 20) // limit upload size to 10MB
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("profile_image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Save the file on disk
	filePath, err := saveFile(file, handler)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Update user profile with the new image
	user.ProfileImage = filePath
	if err := models.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user profile image", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Profile image updated successfully for user %s\n", username)
}

// saveFile saves the uploaded profile image to the uploads directory
func saveFile(file multipart.File, handler *multipart.FileHeader) (string, error) {
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.Mkdir(uploadsDir, os.ModePerm)
	}

	filePath := filepath.Join(uploadsDir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file content to the destination
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filePath, nil
}
