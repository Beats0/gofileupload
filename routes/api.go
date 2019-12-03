package routes

import (
	"github.com/beats0/gofileupload/controllers"
	"github.com/beats0/gofileupload/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.HTML(200, "index.html", gin.H{})
	})
	r.GET("/upload", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.HTML(200, "fileUpload.html", gin.H{})
	})
	r.GET("/register", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.HTML(200, "register.html", gin.H{})
	})
	r.GET("/login", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.HTML(200, "login.html", gin.H{})
	})
	r.GET("/profile", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.HTML(200, "profile.html", gin.H{})
	})

	// api v1 jwt
	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.JwtAuth())
	{
		// 上传文件
		apiV1.POST("/upload", controllers.FormUpload)
		// 多个文件上传
		apiV1.POST("/multiUpload", controllers.MultiUpload)

		apiV1.GET("/user/profile", controllers.GetUserProfile)
		apiV1.GET("/user/fileList", controllers.GetUserUploadByUid)
	}
	// api v2 common
	apiV2 := r.Group("/api/v2")
	{
		apiV2.POST("/user/register", controllers.PostRegister)
		apiV2.POST("/user/login", controllers.PostLogin)
	}
	return r
}
