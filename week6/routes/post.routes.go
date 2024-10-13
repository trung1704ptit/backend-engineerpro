package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/trung/backend-engineerpro/controllers"
	"github.com/trung/backend-engineerpro/middleware"
)

type PostRouteController struct {
	postController controllers.PostController
}

func NewRoutePostController(postController controllers.PostController) PostRouteController {
	return PostRouteController{postController}
}

func (pc *PostRouteController) PostRoute(rg *gin.RouterGroup) {

	router := rg.Group("posts")
	router.POST("", middleware.DeserializeUser(), pc.postController.CreatePost)
	router.GET("", pc.postController.FindPosts)
	router.PUT(":postId", middleware.DeserializeUser(), pc.postController.UpdatePost)
	router.GET(":postId", pc.postController.FindPostById)
	router.DELETE(":postId", middleware.DeserializeUser(), pc.postController.DeletePost)

	router.POST(":postId/like", middleware.DeserializeUser(), pc.postController.ToggleLike)

	comments := router.Group(":postId/comments")
	{
		comments.POST("", middleware.DeserializeUser(), pc.postController.AddComment)
		comments.PUT(":commentId", middleware.DeserializeUser(), pc.postController.UpdateComment)
		comments.DELETE(":commentId", middleware.DeserializeUser(), pc.postController.DeleteComment)
	}
}
