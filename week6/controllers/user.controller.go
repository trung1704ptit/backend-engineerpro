package controllers

import (
	"net/http"
	"time"

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

func (uc *UserController) FollowUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// get user id that will set follow for current user
	userID := ctx.Param("userID")

	// check if user exists
	var userToFollow models.User

	if err := uc.DB.First(&userToFollow, "id = ?", userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User that user wants to follow not found"})
		return
	}

	// check if current user is already following that user
	var follow models.UserFollower
	err := uc.DB.Where("follower_id = ? AND following_id = ?", currentUser.ID, userToFollow.ID).First(&follow).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "You are already following this user"})
		return
	}

	// Add new follow entry in the user_followers table
	follow = models.UserFollower{
		FollowerID:  currentUser.ID,
		FollowingID: userToFollow.ID,
		CreatedAt:   time.Now(),
	}

	if err := uc.DB.Create(&follow).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to follow user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Successfully followed the user"})
}

func (uc *UserController) UnfollowerUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userID := ctx.Param("userID")

	var userToUnfollow models.User
	// check if user exists
	if err := uc.DB.First(&userToUnfollow, "id = ?", userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "User not found"})
		return
	}

	// check if current user is actually following the user
	var follow models.UserFollower
	err := uc.DB.Where("follower_id = ? AND following = ?", currentUser.ID, userToUnfollow.ID).First(&follow).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "You are not following this user"})
		return
	}

	// delete the follow entry from user_followers table
	if err := uc.DB.Delete(&follow).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to unfollow user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Successfully unfollowed the user"})
}
