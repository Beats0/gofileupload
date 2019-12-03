package controllers

import (
	"github.com/beats0/gofileupload/conf"
	myJwt "github.com/beats0/gofileupload/middleware"
	"github.com/beats0/gofileupload/models"
	"github.com/beats0/gofileupload/pkg/e"
	"github.com/beats0/gofileupload/utils"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

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

func PostRegister(ctx *gin.Context) {
	var registerForm RegisterForm
	if err := ctx.ShouldBind(&registerForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(registerForm.Pwd) <= 6 || len(registerForm.Pwd) >= 18 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "密码长度错误，密码长度为6~18",
		})
		return
	}

	uname := ctx.PostForm("uname")
	mail := ctx.PostForm("mail")
	pwd := ctx.PostForm("pwd")

	isMail := utils.VerifyEmail(mail)
	if !isMail {
		ctx.JSON(200, gin.H{
			"code": e.INVALID_MAIL_VERIFY,
			"msg":  e.GetMsg(e.INVALID_MAIL_VERIFY),
		})
		return
	}
	// 如果用户存在
	isResisted := models.IsMailExist(mail)
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
		Pwd:   utils.PwdSaltMd5(pwd),
	}
	err := models.SaveUser(&userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
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
	//ctx.SetCookie("gin_cookie", string(userData.UId), 3600, "/", "localhost", false, true)
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
		"data": userDataTokenRes,
	})
}

func PostLogin(ctx *gin.Context) {
	var LoginForm LoginForm
	if err := ctx.ShouldBind(&LoginForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mail := ctx.PostForm("mail")
	pwd := ctx.PostForm("pwd")

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
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
		},
	}

	token, err := j.CreateToken(claims)
	return token, err
}
