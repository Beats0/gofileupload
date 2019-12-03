package controllers

import (
	"crypto/md5"
	"fmt"
	myJwt "github.com/beats0/gofileupload/middleware"
	"github.com/beats0/gofileupload/models"
	"github.com/beats0/gofileupload/pkg/e"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// 上传文件夹地址
const BASE_NAME = "./static/upload/"

// 返回 url 地址
const RELATIVEPATH_BASE_NAME = "/static/upload/"

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type MultiUploadUrl struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

type UserFile struct {
	Id        uint   `json:"id"`
	UId       uint   `json:"uid"`
	Date      int64  `json:"date"`
	File_name string `json:"file_name"`
	File_size int64  `json:"file_size"`
	Md5       string `json:"md5"`
}

type UserFileListRes struct {
	page  int
	pages int
	list  []*UserFile
}

// 单个文件上传
func FormUpload(ctx *gin.Context) {
	//fh, err := ctx.FormFile("file")
	//checkError(err)
	//ctx.SaveUploadedFile(fh,BASE_NAME + fh.Filename)

	//file, err := fh.Open()
	//defer file.Close()
	//bytes, e := ioutil.ReadAll(file)
	//e = ioutil.WriteFile(BASE_NAME+fh.Filename, bytes, 0666)
	//checkError(e)
	//
	//if e != nil {
	//	ctx.JSON(200, gin.H{
	//		"success": false,
	//	})
	//} else {
	//	ctx.JSON(200, gin.H{
	//		"success": true,
	//	})
	//}

	// single file
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	file, _ := ctx.FormFile("file")
	if file == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR_UPLOAD_SAVE_IMAGE_FAIL,
			"msg":  "上传文件不存在",
		})
		return
	}
	fileHash := GenerateRandomMd5()
	fileExt := strings.ToLower(path.Ext(file.Filename))
	fileType := FileCategory(fileExt)
	savePath := BASE_NAME + fileHash + fileExt
	// TODO: 保存文件分类
	err := ctx.SaveUploadedFile(file, savePath)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  err,
		})
	} else {
		uploadData := models.Upload{
			UId:       claims.UId,
			Date:      time.Now().Unix(),
			File_name: file.Filename,
			File_size: file.Size,
			File_ext:  fileExt,
			File_type: fileType,
			Md5:       fileHash,
		}
		err := models.SaveUploadFile(&uploadData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": e.ERROR_UPLOAD_SAVE_IMAGE_FAIL,
				"msg":  e.GetMsg(e.ERROR_UPLOAD_SAVE_IMAGE_FAIL),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"code": 0,
			"msg":  "0",
			"url":  RELATIVEPATH_BASE_NAME + fileHash + fileExt,
		})
	}
}

// 多个文件上传
func MultiUpload(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	form, err := ctx.MultipartForm()
	checkError(err)
	files := form.File["files"]

	if files == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR_UPLOAD_SAVE_IMAGE_FAIL,
			"msg":  "上传文件不存在",
		})
		return
	}
	var er error
	// TODO: 多文件上传
	var multiUploadUrls []*MultiUploadUrl
	for _, file := range files {
		fileHash := GenerateRandomMd5()
		fileExt := strings.ToLower(path.Ext(file.Filename))
		fileType := FileCategory(fileExt)
		savePath := BASE_NAME + fileHash + fileExt

		er = ctx.SaveUploadedFile(file, savePath)
		multiUploadUrl := MultiUploadUrl{
			Url:   RELATIVEPATH_BASE_NAME + fileHash + fileExt,
			Title: file.Filename,
		}

		uploadData := models.Upload{
			UId:       claims.UId,
			Date:      time.Now().Unix(),
			File_name: file.Filename,
			File_size: file.Size,
			File_ext:  fileExt,
			File_type: fileType,
			Md5:       fileHash,
		}
		err := models.SaveUploadFile(&uploadData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": e.ERROR_UPLOAD_SAVE_IMAGE_FAIL,
				"msg":  e.GetMsg(e.ERROR_UPLOAD_SAVE_IMAGE_FAIL),
			})
			return
		}
		multiUploadUrls = append(multiUploadUrls, &multiUploadUrl)
	}
	if er != nil {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  err,
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
		"url":  multiUploadUrls,
	})
}

func checkError(error error) {
	if error != nil {
		fmt.Println(error)
	}
	return
}

// 获取用户上传文件列表
func GetUserUploadByUid(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)

	var (
		data  = make(map[string]interface{})
		page  int
		limit int
	)
	if ctx.Query("page") == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(ctx.Query("page"))
	}
	if ctx.Query("limit") == "" {
		limit = 25
	} else {
		limit, _ = strconv.Atoi(ctx.Query("limit"))
	}

	userFiles, err := models.GetUserUploadByUid(claims.UId, page, limit)
	total, err := models.GetUserUploadTotalByUid(claims.UId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "获取用户文件列表失败",
		})
		return
	}
	pages := math.Ceil(float64(total / limit))
	data["lists"] = userFiles
	data["page"] = page
	data["pages"] = pages
	data["total"] = total
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
		"data": data,
	})
}

// 返回文件hash
func GetFileHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	fileHash := fmt.Sprintf("%x", h.Sum(nil))
	rnameErr := os.Rename(filePath, BASE_NAME+fileHash+".png")
	if rnameErr != nil {
		return "", rnameErr
	}
	return fileHash, nil
}

// 生成随机md5 32
func GenerateRandomMd5(n ...int) string {
	noRandomCharacters := 32
	if len(n) > 0 {
		noRandomCharacters = n[0]
	}
	randString := GenerateRandomStr(noRandomCharacters)
	Md5Inst := md5.New()
	Md5Inst.Write([]byte(randString))
	bs := Md5Inst.Sum([]byte(""))
	return fmt.Sprintf("%x", bs)
}

// 生成随机字符串 32
func GenerateRandomStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}

// 文件分类
func FileCategory(fileExt string) string {
	fileExt = strings.ToLower(fileExt)
	var category string
	switch fileExt {
	case ".png", ".jpg", ".jpeg", ".bmp", ".gif", ".webp":
		category = "image"
	case ".mp3", ".ogg":
		category = "audio"
	case ".mp4":
		category = "media"
	}
	return category
}
