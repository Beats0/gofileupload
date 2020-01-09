package controllers

import (
	"crypto/md5"
	"fmt"
	"github.com/beats0/gofileupload/conf"
	myJwt "github.com/beats0/gofileupload/middleware"
	"github.com/beats0/gofileupload/models"
	"github.com/beats0/gofileupload/pkg/e"
	"github.com/beats0/gofileupload/utils"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
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
const BASE_PATH = "./static/upload/"

// 返回 url 地址
const RELATIVEPATH_BASE_PATH = "/static/upload/"

var characterRunes = []rune("abcdefghijklmnopqueryFormrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

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

type QueryForm struct {
	Order    string `form:"order,default='date'" json:"order,default='date'"`
	DirId    uint   `form:"dirId,default=0" json:"dirId,default=0"`
	Limit    int    `form:"limit,default=50" json:"limit,default=50" binding:"max=100,min=50"`
	Desc     int    `form:"desc,default=1" json:"desc,default=1"`
	Page     int    `form:"page,default=1" json:"page,default=1"`
	FileType string `form:"fileType" json:"fileTye"`
	Q        string `form:"q" json:"q"`
}

type UploadForm struct {
	Index      int    `form:"index" json:"index"`
	FileMD5    string `form:"fileMD5" json:"fileMD5"`
	DirId      uint   `form:"dirId,default=0" json:"dirId,default=0"`
	UploadType string `form:"uploadType" json:"uploadType"`
	FileName   string `form:"fileName" json:"fileName"`
	FileSize   int64  `form:"fileSize" json:"fileSize"`
}

type CreateDirForm struct {
	DirId   uint   `form:"dirId" json:"dirId"`
	DirName string `form:"dirName" json:"dirName" binding:"required"`
}

type RenameForm struct {
	FsId     uint   `form:"fsId" json:"fsId"`
	DirId    uint   `form:"dirId" json:"dirId"`
	FileName string `form:"fileName" json:"fileName" binding:"required"`
}

// 单个文件上传
func FormUpload(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	file, _ := ctx.FormFile("file")
	if file == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR_UPLOAD_SAVE_IMAGE_FAIL,
			"msg":  "上传文件为空",
		})
		return
	}
	fileHash := GenerateRandomMd5()
	fileExt := strings.ToLower(path.Ext(file.Filename))
	fileType := FileCategory(fileExt)
	savePath := BASE_PATH + fileHash + fileExt
	// TODO: 保存文件分类
	err := ctx.SaveUploadedFile(file, savePath)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  err,
		})
	} else {
		uploadData := models.Upload{
			UId:           claims.UId,
			Date:          time.Now().Unix(),
			File_name:     file.Filename,
			File_size:     file.Size,
			File_ext:      fileExt,
			File_type:     fileType,
			Md5:           fileHash,
			Last_modified: time.Now().Unix(),
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
			"url":  RELATIVEPATH_BASE_PATH + fileHash + fileExt,
		})
	}
}

// 多个文件上传
func MultiUpload(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  err,
		})
		return
	}
	files := form.File["files"]

	if files == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR_UPLOAD_SAVE_IMAGE_FAIL,
			"msg":  "上传文件为空",
		})
		return
	}
	var er error
	var multiUploadUrls []*MultiUploadUrl
	for _, file := range files {
		fileHash := GenerateRandomMd5()
		fileExt := strings.ToLower(path.Ext(file.Filename))
		fileType := FileCategory(fileExt)
		savePath := BASE_PATH + fileHash + fileExt

		er = ctx.SaveUploadedFile(file, savePath)
		multiUploadUrl := MultiUploadUrl{
			Url:   RELATIVEPATH_BASE_PATH + fileHash + fileExt,
			Title: file.Filename,
		}

		uploadData := models.Upload{
			UId:           claims.UId,
			Date:          time.Now().Unix(),
			File_name:     file.Filename,
			File_size:     file.Size,
			File_ext:      fileExt,
			File_type:     fileType,
			Md5:           fileHash,
			Last_modified: time.Now().Unix(),
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

// 分片上传
func SliceUpload(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	var uploadForm UploadForm
	if err := ctx.ShouldBind(&uploadForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	file, _ := ctx.FormFile("file")
	var (
		uid        = claims.UId
		dirId      = uploadForm.DirId
		fileIndex  = uploadForm.Index
		fileMD5    = uploadForm.FileMD5
		uploadType = uploadForm.UploadType
		fileName   = uploadForm.FileName
		fileSize   = uploadForm.FileSize
	)
	if file == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR_UPLOAD_SAVE_IMAGE_FAIL,
			"msg":  "上传文件不存在",
		})
		return
	}
	// 保存到对应fileMD5文件夹中
	saveDir := BASE_PATH + fileMD5
	// 分片上传保存
	if uploadType == "slice" {
		_, errDir := utils.CheckDirPath(saveDir)
		if errDir != nil {
			ctx.JSON(500, gin.H{
				"code": 500,
				"msg":  errDir,
			})
			return
		}
		savePath := saveDir + "/" + fileMD5 + "." + strconv.Itoa(fileIndex) + ".part"
		err := ctx.SaveUploadedFile(file, savePath)
		if err != nil {
			ctx.JSON(500, gin.H{
				"code": 500,
				"msg":  err,
			})
		} else {
			ctx.JSON(200, gin.H{
				"code":  0,
				"msg":   0,
				"index": fileIndex,
			})
		}
	} else if uploadType == "merge" {
		partFiles, dirErr := utils.GetFilePart(saveDir)
		if dirErr != nil {
			// 合并分片文件
			ctx.JSON(500, gin.H{
				"code": 500,
				"msg":  "文件夹不存在",
			})
			return
		}

		fileExt := strings.ToLower(path.Ext(fileName))
		mergeFilePath := BASE_PATH + fileMD5 + fileExt
		fileType := FileCategory(fileExt)
		// 合并分片文件并合并
		f, _ := os.OpenFile(mergeFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		for _, v := range partFiles {
			contents, _ := ioutil.ReadFile(v)
			f.Write(contents)
			// NOTE:
			// 是否需要删除分片文件
			// os.Remove(v)
			// 是否需要删除分片文件夹
			// os.RemoveAll(BASE_PATH + fileMD5)
		}
		defer f.Close()
		// 检查是否有重复文件名
		existTotal, err := models.IsItemExist(dirId, uid, fileName)
		if err == nil && existTotal != 0 {
			fileName += "(副本)"
		}
		// 保存数据库
		uploadData := models.Upload{
			UId:           uid,
			Date:          time.Now().Unix(),
			File_name:     fileName,
			File_size:     fileSize,
			File_ext:      fileExt,
			File_type:     fileType,
			Md5:           fileMD5,
			Parent:        dirId,
			Last_modified: time.Now().Unix(),
		}
		err = models.SaveUploadFile(&uploadData)
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
			"url":  RELATIVEPATH_BASE_PATH + fileMD5 + fileExt,
		})
	}
}

// 检查分片进度
func CheckFile(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	UId := claims.UId
	fileMD5 := ctx.Query("fileMD5")
	_, err := models.CheckUserFileMd5(fileMD5, UId)
	// 没有记录，返回空
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 0,
			"msg":  "0",
		})
		return
	}
	targetDir := BASE_PATH + fileMD5
	// part 文件最大索引值
	partFiles, _ := utils.GetFilePart(targetDir)
	ctx.JSON(200, gin.H{
		"code":  0,
		"msg":   "0",
		"index": len(partFiles) - 1,
	})
}

// 创建文件夹
func CreateFolder(ctx *gin.Context) {
	var createDirForm CreateDirForm
	if err := ctx.ShouldBind(&createDirForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	var (
		dirId   = createDirForm.DirId
		dirName = createDirForm.DirName
		uid     = claims.UId
	)

	uploadData := models.Upload{
		Parent:        dirId,
		UId:           uid,
		Date:          time.Now().Unix(),
		File_name:     dirName,
		Is_dir:        1,
		Last_modified: time.Now().Unix(),
	}
	//是否含有该文件夹
	existTotal, err := models.IsItemExist(dirId, uid, dirName)
	if err == nil && existTotal != 0 {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "新建文件夹失败, 已包含该文件夹",
		})
		return
	}
	// 数据库保存创建文件夹数据
	uploadDataRes, err := models.CreateFolder(uploadData)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "新建文件夹失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
		"data": uploadDataRes,
	})
}

// 获取用户上传文件列表
func FileList(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	var queryForm QueryForm
	if err := ctx.ShouldBind(&queryForm); err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	var (
		uid   = claims.UId
		dirId = queryForm.DirId
		order = queryForm.Order
		page  = queryForm.Page
		limit = queryForm.Limit
		data  = make(map[string]interface{})
		desc  string
	)
	if queryForm.Desc == 1 {
		desc = "ASC"
	} else {
		desc = "DESC"
	}
	userFiles, err := models.GetUserUploadByUid(uid, dirId, order, desc, page, limit)
	total, err := models.GetUserUploadTotalByUid(dirId, uid)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取用户文件列表失败",
		})
		return
	}
	pages := math.Ceil(float64(total) / float64(limit))
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

// 文件分类
func FileListByFileType(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	var queryForm QueryForm
	if err := ctx.ShouldBind(&queryForm); err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	var (
		uid      = claims.UId
		fileType = queryForm.FileType
		order    = queryForm.Order
		page     = queryForm.Page
		limit    = queryForm.Limit
		data     = make(map[string]interface{})
		desc     string
	)
	if queryForm.Desc == 1 {
		desc = "ASC"
	} else {
		desc = "DESC"
	}
	userFiles, err := models.GetUserFileTypesByUid(uid, fileType, order, desc, page, limit)
	total, err := models.GetUserFileTypeTotal(fileType, uid)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取用户文件列表失败",
		})
		return
	}
	pages := math.Ceil(float64(total) / float64(limit))
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

// 文件搜索
func Search(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	var queryForm QueryForm
	if err := ctx.ShouldBind(&queryForm); err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	var (
		uid   = claims.UId
		q     = queryForm.Q
		order = queryForm.Order
		page  = queryForm.Page
		limit = queryForm.Limit
		data  = make(map[string]interface{})
		desc  string
	)
	if queryForm.Desc == 1 {
		desc = "ASC"
	} else {
		desc = "DESC"
	}
	userFiles, err := models.Search(uid, q, order, desc, page, limit)
	total, err := models.SearchTotal(uid, q)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取用户文件列表失败",
		})
		return
	}
	pages := math.Ceil(float64(total) / float64(limit))
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

// 重命名
func Rename(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*myJwt.CustomClaims)
	var renameForm RenameForm
	if err := ctx.ShouldBind(&renameForm); err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	var (
		uid      = claims.UId
		fsId     = renameForm.FsId
		fileName = renameForm.FileName
	)

	userFiles, err := models.GetUserFileByDirAndFileName(uid, fsId, fileName)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 500,
			"msg":  "获取用户文件列表失败",
		})
		return
	}

	// 同一文件夹下，名字相同且不为目标文件，改为副本
	for _, v := range userFiles {
		if v.File_name == fileName && v.ID != fsId {
			fileName += "(副本)"
		}
	}

	if err := models.Rename(fsId, uid, fileName); err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
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
	rnameErr := os.Rename(filePath, BASE_PATH+fileHash+".png")
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
	case ".mp3", ".wav", ".ogg":
		category = "audio"
	case ".mp4", ".rm", ".rmvb", ".mpg", ".avi", ".mov", ".wmv", ".flv":
		category = "media"
	case ".txt", ".doc", ".docx", ".pdf", ".xls", ".xlsx", ".ppt", ".pptx":
		category = "doc"
	case ".zip", ".rar", ".7z", ".tar", ".gz":
		category = "zip"
	case ".torrent":
		category = "torrent"
	default:
		category = "other"
	}
	return category
}

// 获取服务器磁盘配置
func ServerStatic(ctx *gin.Context) {
	var (
		data = make(map[string]interface{})
	)
	totalSize, err := utils.DirSize("static")
	if err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "获取static大小失败",
		})
		return
	}
	data["totalSize"] = totalSize
	data["maxDiskSize"] = conf.ServerSetting.MaxDiskSize
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "0",
		"data": data,
	})
}
