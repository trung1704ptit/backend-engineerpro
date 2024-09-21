package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trung/backend-engineerpro/models"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		ID:           currentUser.ID,
		Name:         currentUser.Name,
		Age:          currentUser.Age,
		Email:        currentUser.Email,
		ProfileImage: currentUser.ProfileImage,
		Role:         currentUser.Role,
		Provider:     currentUser.Provider,
		CreatedAt:    currentUser.CreatedAt,
		UpdatedAt:    currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}
