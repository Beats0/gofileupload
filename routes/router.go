package routes

import (
	"github.com/beats0/gofileupload/controllers"
	"github.com/beats0/gofileupload/middleware"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 以 assets 作为静态static, 可以直接访问 static/assets, 不能访问 static/upload
	r.Use(static.Serve("/", static.LocalFile("./static/assets", false)))
	r.Use(middleware.CORSMiddleware())
	r.LoadHTMLGlob("templates/*")

	// go vue-router rewrite
	r.GET("/", RouterRewrite)
	r.GET("/login", RouterRewrite)
	r.GET("/register", RouterRewrite)
	r.GET("/disk", RouterRewrite)
	r.GET("/share", RouterRewrite)
	r.GET("/play", RouterRewrite)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "404", "message": "404 not found2", "requestUrl": c.Request.URL})
	})

	// uri
	// 图片 filter
	r.GET("/image/:imagename", controllers.ImageFilter)
	// 下载文件
	// TODO: 检验文件权限
	r.GET("/download/:filename", controllers.DownloadFile)
	// 视频
	r.GET("/video/:v", middleware.CheckSecretKey, controllers.Video)

	// api v1 auth jwt
	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.JwtAuth())
	{
		// 个人信息
		apiV1.GET("/user/profile", controllers.GetUserProfile)
		// 个人文件列表
		apiV1.GET("/user/fileList", controllers.FileList)
		// 文件分类列表
		apiV1.GET("/user/fileListByFileType", controllers.FileListByFileType)
		// 上传文件
		apiV1.POST("/upload", controllers.FormUpload)
		// 新建文件夹
		apiV1.POST("/createFolder", controllers.CreateFolder)
		// 多个文件上传
		apiV1.POST("/multiUpload", controllers.MultiUpload)
		// 分片上传
		apiV1.POST("/sliceUpload", controllers.SliceUpload)
		// 检查分片进度
		apiV1.GET("/checkFile", controllers.CheckFile)
		// 删除文件
		apiV1.POST("/deleteFile", controllers.DeleteFile)
		// 删除文件多个
		apiV1.POST("/deleteFileArr", controllers.DeleteFileArr)
		// 查看文件所在位置
		apiV1.POST("/findFilePath", controllers.FindFilePath)
		// 文件搜索
		apiV1.POST("/search", controllers.Search)
		// 文件命名
		apiV1.POST("/rename", controllers.Rename)
		// 视频信息
		apiV1.GET("/videoInfo", controllers.VideoInfo)
		// 服务器信息
		apiV1.GET("/serverStatic", controllers.ServerStatic)
	}
	// api v2 common
	apiV2 := r.Group("/api/v2")
	{
		// 登录
		apiV2.POST("/user/login", controllers.PostLogin)
		// 注册
		apiV2.POST("/user/register", controllers.PostRegister)
	}
	return r
}

// vue-router rewrite
func RouterRewrite(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.HTML(200, "index.html", gin.H{})
}
