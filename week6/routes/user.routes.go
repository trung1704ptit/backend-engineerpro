package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/trung/backend-engineerpro/controllers"
	"github.com/trung/backend-engineerpro/middleware"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.GET("/profile", middleware.DeserializeUser(), uc.userController.UserProfile)
	router.PUT("/profile", middleware.DeserializeUser(), uc.userController.UpdateUserProfile)
	router.POST("/follow/:userID", middleware.DeserializeUser(), uc.userController.FollowUser)
	router.DELETE("/unfollow/:userID", middleware.DeserializeUser(), uc.userController.UnfollowerUser)
}
