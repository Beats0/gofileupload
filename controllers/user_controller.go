package controllers

import (
	"bufio"
	"fmt"
	"github.com/beats0/gofileupload/conf"
	myJwt "github.com/beats0/gofileupload/middleware"
	"github.com/beats0/gofileupload/models"
	"github.com/beats0/gofileupload/pkg/e"
	"github.com/beats0/gofileupload/utils"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	_ "io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const DOWNLOADS_PATH = "static/upload"

type LoginResult struct {
	Token string `json:"token"`
}

type UserData struct {
	UId    uint   `json:"uid"`
	Uname  string `json:"uname"`
	Avatar string `json:"avatar"`
}

type UserDataTokenRes struct {
	UId    uint   `json:"uid"`
	Uname  string `json:"uname"`
	Avatar string `json:"avatar"`
	Token  string `json:"token"`
}

type LoginForm struct {
	Mail string `form:"mail" json:"mail" binding:"required"`
	Pwd  string `form:"pwd" json:"pwd" binding:"required"`
}

type RegisterForm struct {
	Uname string `form:"uname" json:"uname" binding:"required"`
	Mail  string `form:"mail" json:"mail" binding:"required"`
	Pwd   string `form:"pwd" json:"pwd" binding:"required"`
}

type DeleteFileArrForm struct {
	DeleteFileArr []uint `form:"deleteFileArr" json:"deleteFileArr" binding:"required"`
}

type DeleteForm struct {
	FsId uint `form:"fs_id" json:"fs_id" binding:"required"`
}

func PostRegister(ctx *gin.Context) {
	var registerForm RegisterForm
	if err := ctx.ShouldBind(&registerForm); err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	if len(registerForm.Pwd) <= 6 || len(registerForm.Pwd) >= 18 {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "密码长度错误，密码长度为6~18",
		})
		return
	}

	var (
		uname = registerForm.Uname
		mail  = registerForm.Mail
		pwd   = registerForm.Pwd
	)

	isMail := utils.VerifyEmail(mail)
	if !isMail {
		ctx.JSON(200, gin.H{
			"code": e.INVALID_MAIL_VERIFY,
			"msg":  e.GetMsg(e.INVALID_MAIL_VERIFY),
		})
		return
	}
	// 如果用户存在
	isResisted, _ := models.IsMailExist(mail)
	if isResisted {
		ctx.JSON(200, gin.H{
			"code": e.ERROR_USERMAIL_EXIST,
			"msg":  e.GetMsg(e.ERROR_USERMAIL_EXIST),
		})
		return
	}

	userData := models.User{
		Uname: uname,
		Date:  time.Now().Unix(),
		Mail:  mail,
		// 密码加密
		Pwd: utils.PwdSaltMd5(pwd),
	}
	err := models.SaveUser(&userData)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": e.ERROR_ADDUSER_FAIL,
			"msg":  e.GetMsg(e.ERROR_ADDUSER_FAIL),
		})
		return
	}
	userQueryData, err := models.GetUserByMail(mail)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": e.ERROR_GET_USER_FAIL,
			"msg":  e.GetMsg(e.ERROR_GET_USER_FAIL),
		})
		return
	}
	token, err := generateToken(userQueryData.UId)
	userDataTokenRes := UserDataTokenRes{
		UId:    userQueryData.UId,
		Uname:  userQueryData.Uname,
		Avatar: userQueryData.Avatar,
		Token:  token,
	}
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
		"data": userDataTokenRes,
	})
}

func PostLogin(ctx *gin.Context) {
	var loginForm LoginForm
	if err := ctx.ShouldBind(&loginForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	var (
		mail = loginForm.Mail
		pwd  = loginForm.Pwd
	)
	userData, err := models.UserLogin(mail, pwd)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": e.ERROR_GET_USER_FAIL,
			"msg":  e.GetMsg(e.ERROR_GET_USER_FAIL),
		})
		return
	}
	if userData == nil {
		ctx.JSON(200, gin.H{
			"code": e.ERROR_LOGIN_PWD_FAIL,
			"msg":  e.GetMsg(e.ERROR_LOGIN_PWD_FAIL),
		})
		return
	}
	token, err := generateToken(userData.UId)
	userDataTokenRes := UserDataTokenRes{
		UId:    userData.UId,
		Uname:  userData.Uname,
		Avatar: userData.Avatar,
		Token:  token,
	}
	//ctx.SetCookie("gin_cookie", string(userData.UId), 3600, "/", "localhost", false, true)
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
		"data": userDataTokenRes,
	})
}

func GetUserProfile(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	if claims != nil {
		userData, err := models.GetUserByUid(claims.UId)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": 0,
				"msg":    "获取用户信息错误",
			})
			return
		}
		userResData := UserData{
			UId:    userData.UId,
			Uname:  userData.Uname,
			Avatar: userData.Avatar,
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "0",
			"data":   userResData,
		})
	}
}

// 生成令牌
func generateToken(uid uint) (string, error) {
	j := &myJwt.JWT{
		[]byte(conf.AppSetting.JwtSecret),
	}
	claims := myJwt.CustomClaims{
		UId: uid,
		StandardClaims: jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),      // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600*24*7), // 过期时间 7天
		},
	}

	token, err := j.CreateToken(claims)
	return token, err
}

// 删除文件
func DeleteFile(ctx *gin.Context) {
	var deleteForm DeleteForm
	if err := ctx.ShouldBind(&deleteForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	var (
		uid  = claims.UId
		fsId = deleteForm.FsId
	)

	userFile, err := models.CheckUserFile(fsId, uid)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 0,
			"msg":  "删除失败， 项目不存在",
		})
		return
	}
	if err := models.DeleteUserFile(fsId, uid); err != nil {
		ctx.JSON(200, gin.H{
			"code": 0,
			"msg":  "删除失败， 项目不存在",
		})
		return
	}

	md5 := userFile.Md5
	ext := userFile.File_ext
	// 单个删除模式下，是否需要删除文件夹，默认是
	targetPath := filepath.Join(DOWNLOADS_PATH, md5+ext)
	targetDirPath := filepath.Join(DOWNLOADS_PATH, md5)
	fmt.Println("删除文件", userFile.File_name, targetPath)
	fmt.Println("删除文件夹", targetDirPath)
	if err := os.Remove(targetPath); err != nil {
		ctx.JSON(200, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("%s删除文件失败,err:%v", userFile.File_name, err),
		})
		return
	}
	if err := os.RemoveAll(targetDirPath); err != nil {
		ctx.JSON(200, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("%s删除文件夹失败,err:%v", userFile.File_name, err),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
	})
}

// 批量删除多个文件
func DeleteFileArr(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	var deleteFileArrForm DeleteFileArrForm
	if err := ctx.ShouldBind(&deleteFileArrForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	var (
		uid           = claims.UId
		deleteFileArr = deleteFileArrForm.DeleteFileArr
	)

	for _, fsId := range deleteFileArr {
		userFile, _ := models.CheckUserFile(fsId, uid)
		_ = models.DeleteUserFile(fsId, uid)
		md5 := userFile.Md5
		ext := userFile.File_ext
		// 是否需要删除文件夹，默认否
		targetPath := filepath.Join(DOWNLOADS_PATH, md5+ext)
		targetDirPath := filepath.Join(DOWNLOADS_PATH, md5)
		fmt.Println("删除文件", userFile.File_name, targetPath)
		fmt.Println("删除文件夹", targetDirPath)
		// 批量情况下, 是否需要删除文件夹，默认否
		//if err := os.Remove(targetPath); err != nil {
		//	ctx.JSON(200, gin.H{
		//		"code": 0,
		//		"msg":  fmt.Sprintf("%s删除文件失败,err:%v", userFile.File_name, err),
		//	})
		//	return
		//}
		//if err := os.RemoveAll(targetDirPath); err != nil {
		//	ctx.JSON(200, gin.H{
		//		"code": 0,
		//		"msg":  fmt.Sprintf("%s删除文件夹失败,err:%v", userFile.File_name, err),
		//	})
		//	return
		//}
	}
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
	})
}

// 查看文件所在目录
func FindFilePath(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	var deleteForm DeleteForm
	if err := ctx.ShouldBind(&deleteForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	var (
		uid  = claims.UId
		fsId = deleteForm.FsId
		data = make(map[string]interface{})
	)
	parent, err := models.FindFilePath(fsId, uid)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	data["dirId"] = parent
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
		"data": data,
	})
}

// 图片filter
// TODO: 访问权限
func ImageFilter(ctx *gin.Context) {
	var (
		imagename = ctx.Param("imagename")
		w, _      = strconv.Atoi(ctx.Query("w"))
		h, _      = strconv.Atoi(ctx.Query("h"))
	)
	if imagename == "" {
		ctx.JSON(404, gin.H{
			"code": 404,
			"msg":  "Param error",
		})
		return
	}
	targetPath := filepath.Join(DOWNLOADS_PATH, imagename)
	if !utils.FileExit(targetPath) {
		ctx.JSON(404, gin.H{
			"code": 404,
			"msg":  "Image not found",
		})
		return
	}

	// 没有宽高参数加载原图像
	if w == 0 && h == 0 {
		ctx.Header("Content-Type", "image/jpeg")
		ctx.File(targetPath)
		return
	}

	var (
		re      = regexp.MustCompile(`[a-f0-9]{32}`)
		md5     = re.FindString(targetPath)
		fileExt = path.Ext(targetPath)
	)

	// 裁剪图像路径
	cutPath := fmt.Sprintf("%s_%d_%d%s", md5, w, h, fileExt)
	// 01378a9d30a32b34bd8a1dfd1f99c04a_200_200.jpg
	cutTargetPath := filepath.Join(DOWNLOADS_PATH, cutPath)
	// 判断是否存在裁剪图像
	// 存在裁剪图片直接返回
	if utils.FileExit(cutTargetPath) {
		ctx.Header("Content-Type", "image/jpeg")
		ctx.File(cutTargetPath)
		return
	}

	// 不存在重新裁剪图片
	file, _ := os.Open(targetPath)
	defer file.Close()

	// 图片解码
	bufFile := bufio.NewReader(file)
	img, imgtype, _ := image.Decode(bufFile)
	// 要裁剪的宽高不能大于自身的宽高
	Rwidth := img.Bounds().Max.X
	if w > Rwidth {
		w = Rwidth
	}
	Rheight := img.Bounds().Max.Y
	if h > Rheight {
		h = Rheight
	}

	// 忽略gif
	if imgtype == "gif" || (w == Rwidth && h == Rheight) {
		ctx.Header("Content-Type", "image/gif")
		ctx.File(targetPath)
		return
	}

	// 进行裁剪
	reimg := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	// 裁剪的存储
	out, err := os.Create(cutTargetPath)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "Save cut image error",
		})
		return
	}
	defer out.Close()

	if imgtype == "jpeg" || imgtype == "jpg" {
		// 保存裁剪的图片
		_ = jpeg.Encode(out, reimg, nil)
		// 向浏览器输出
		_ = jpeg.Encode(ctx.Writer, reimg, nil)
	} else if imgtype == "png" {
		// 保存裁剪的图片
		_ = png.Encode(out, reimg)
		// 向浏览器输出
		_ = png.Encode(ctx.Writer, reimg)
	}
}

// 下载文件
// TODO: 访问权限
func DownloadFile(ctx *gin.Context) {
	fileName := ctx.Param("filename")
	if fileName == "" {
		ctx.JSON(404, gin.H{
			"code": 404,
			"msg":  "Param error",
		})
		return
	}
	targetPath := filepath.Join(DOWNLOADS_PATH, fileName)
	if !utils.FileExit(targetPath) {
		ctx.JSON(404, gin.H{
			"code": 404,
			"msg":  "Image not found",
		})
		return
	}
	// 设置header
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(targetPath)
}

// 视频播放
func Video(ctx *gin.Context) {
	var (
		videoMd5 = ctx.Param("v")
	)
	if videoMd5 == "" {
		ctx.JSON(404, gin.H{
			"code": 404,
			"msg":  "Param error",
		})
		return
	}
	sFile, err := models.GetFileByMd5(videoMd5)
	if err != nil {
		ctx.JSON(404, gin.H{
			"code": 404,
			"msg":  "Video not found",
		})
		return
	}
	targetPath := filepath.Join(DOWNLOADS_PATH, videoMd5+sFile.File_ext)
	if !utils.FileExit(targetPath) {
		ctx.JSON(404, gin.H{
			"code": 404,
			"msg":  "Video not found",
		})
		return
	}
	ctx.Header("Content-Type", "video/mp4")
	ctx.File(targetPath)
}

// 视频信息
func VideoInfo(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	var (
		videoMd5 = ctx.Query("v")
		uid      = claims.UId
	)
	userFile, err := models.GetVideoInfo(uid, videoMd5)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 400,
			"msg":  "项目不存在",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "",
		"data": userFile,
	})
}
