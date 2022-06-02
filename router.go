package main

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/test/", controller.Test)
	// user api

	userApi := apiRouter.Group("/user")
	//userApi.Use(sessions.Sessions("mysession", global.DY_SESSION_STORE))
	userApi.POST("/register/", controller.Register)

	// session middleware

	userApi.POST("/login/", controller.Login)

	////session login
	//userApi.GET("", controller.UserInfo)
	//
	//apiRouter.POST("/publish/action/", controller.Publish)
	//apiRouter.GET("/publish/list/", controller.PublishList)

	//jwt logic
	userApi.GET("", utils.JWTAuthMiddleware(), controller.UserInfo)

	apiRouter.POST("/publish/action/", utils.JWTAuthMiddleware(), controller.Publish)
	apiRouter.GET("/publish/list/", utils.JWTAuthMiddleware(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", utils.JWTAuthMiddleware(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", utils.JWTAuthMiddleware(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}
