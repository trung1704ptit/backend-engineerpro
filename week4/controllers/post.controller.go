package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/trung/backend-engineerpro/models"
	"gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

func NewPostController(DB *gorm.DB) PostController {
	return PostController{DB}
}

func (pc *PostController) CreatePost(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreatePostRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newPost := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		UserID:    currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := pc.DB.Create(&newPost)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Post with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPost})
}

func (pc *PostController) UpdatePost(ctx *gin.Context) {
	postId := ctx.Param("postId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdatePost
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedPost models.Post
	result := pc.DB.Where("id = ? AND user_id = ?", postId, currentUser.ID).First(&updatedPost)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}
	now := time.Now()
	postToUpdate := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Image:     payload.Image,
		UserID:    currentUser.ID,
		CreatedAt: updatedPost.CreatedAt,
		UpdatedAt: now,
	}

	pc.DB.Model(&updatedPost).Updates(postToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedPost})
}

func (pc *PostController) FindPostById(ctx *gin.Context) {
	postId := ctx.Param("postId")

	var post models.Post
	result := pc.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

func (pc *PostController) FindPosts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var posts []models.Post
	results := pc.DB.Preload("Likes").Preload("Comments").Limit(intLimit).Offset(offset).Find(&posts)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(posts), "data": posts})
}

func (pc *PostController) DeletePost(ctx *gin.Context) {
	postId := ctx.Param("postId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	result := pc.DB.Where("id = ? AND user_id = ?", postId, currentUser.ID).Delete(&models.Post{})

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (pc *PostController) ToggleLike(ctx *gin.Context) {
	var like models.Like
	postIdStr := ctx.Param("postId")
	currentUser := ctx.MustGet("currentUser").(models.User)
	postId, err := uuid.Parse(postIdStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid post ID format"})
		return
	}

	// check if user already like this page
	result := pc.DB.Where("post_id = ? AND user_id = ?", postId, currentUser.ID).First(&like)

	if result.RowsAffected == 0 {
		// Case 1: The user hasn't liked the post yet, so create new Like
		now := time.Now()
		newLike := models.Like{
			PostID:    postId,
			UserID:    currentUser.ID,
			CreateAt:  now,
			UpdatedAt: now,
		}

		if err := pc.DB.Create(&newLike).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to like the post"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Post like update successfully"})
	} else {
		// Case 2: User already liked the post
		if err := pc.DB.Delete(&like).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to dislike the post"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Post like update successfully"})
	}
}

func (pc *PostController) AddComment(ctx *gin.Context) {
	postIdStr := ctx.Param("postId")
	currentUser := ctx.MustGet("currentUser").(models.User)
	postId, err := uuid.Parse(postIdStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid post ID format"})
		return
	}
	now := time.Now()
	var payload *models.CreateComment

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newComment := models.Comment{
		PostID:    postId,
		UserID:    currentUser.ID,
		Content:   payload.Content,
		CreateAt:  now,
		UpdatedAt: now,
	}

	if err := pc.DB.Create(&newComment).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Could not add new comment"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newComment})
}

func (pc *PostController) UpdateComment(ctx *gin.Context) {
	postIdStr := ctx.Param("postId")
	commentId := ctx.Param("commentId")
	currentUser := ctx.MustGet("currentUser").(models.User)
	now := time.Now()

	var payload *models.UpdateComment
	var updatedComment models.Comment
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Parse the postId from string to UUID
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid post ID format"})
		return
	}

	result := pc.DB.Where("id = ? AND post_id = ? AND user_id = ?", commentId, postId, currentUser.ID).First(&updatedComment)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Comment not exists"})
		return
	}

	commentToUpdate := models.Comment{
		PostID:    postId,
		UserID:    currentUser.ID,
		UpdatedAt: now,
		Content:   payload.Content,
	}

	pc.DB.Where("id = ? AND post_id = ? AND user_id = ?", commentId, postId, currentUser.ID).Model(&updatedComment).Updates(commentToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedComment})
}

func (pc *PostController) DeleteComment(ctx *gin.Context) {
	postIdStr := ctx.Param("postId")
	commentId := ctx.Param("commentId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	// Parse the postId from string to UUID
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid post ID format"})
		return
	}

	// Attempt to delete the comment with the matching commentId, postId, and userId
	result := pc.DB.Where("id = ? AND post_id = ? AND user_id = ?", commentId, postId, currentUser.ID).Delete(&models.Comment{})

	// Check if any row was affected (meaning comment was found and deleted)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Comment not found or not authorized to delete"})
		return
	}

	// Handle any potential database errors
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something went wrong, cannot delete comment"})
		return
	}

	// Successful deletion, returning no content
	ctx.JSON(http.StatusNoContent, nil)
}
