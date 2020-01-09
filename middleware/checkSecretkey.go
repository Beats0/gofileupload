package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/beats0/gofileupload/conf"
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
	"time"
)

func CheckSecretKey(ctx *gin.Context) {
	var (
		req     = ctx.Request
		q       = req.URL.Query()
		expTime = q.Get("exptime")
		secret  = q.Get("secret")
		reqHost = req.Host
		referer = req.Header.Get("Referer")
	)

	// release 模式下对 reqHost 和 referer 验证
	if conf.ServerSetting.RunMode == "release" {
		uri, err := url.Parse(referer)
		if err != nil {
			ctx.JSON(403, gin.H{
				"code": 403,
				"msg":  "parse referer error",
			})
			ctx.Abort()
			return
		}
		serverUrl, err := url.Parse(conf.AppSetting.Url)
		if err != nil {
			ctx.JSON(403, gin.H{
				"code": 403,
				"msg":  "App conf Url parse error",
			})
			ctx.Abort()
			return
		}
		if reqHost != serverUrl.Host || uri.Host != serverUrl.Host {
			ctx.JSON(403, gin.H{
				"code": 403,
				"msg":  "Host or referer error",
			})
			ctx.Abort()
			return
		}
	}
	reqUrl := fmt.Sprintf("http://%s%s", reqHost, req.URL.Path)

	key := conf.AppSetting.CheckSecretKey
	err := check(secret, expTime, reqUrl, key)

	if err != nil {
		ctx.JSON(403, gin.H{
			"code": 403,
			"msg":  "token expired: " + err.Error(),
		})
		ctx.Abort()
		return
	}
	ctx.Next()
}

func check(secret string, expTime string, url string, key string) error {
	if secret == "" || expTime == "" {
		return errors.New("secret or time is empty")
	}

	expTimeInt, err := strconv.ParseInt(expTime, 10, 64)
	// 过期时间 12 小时
	sExpTime := time.Now().Add(-12 * time.Hour).Unix()
	if err != nil || sExpTime > expTimeInt {
		return errors.New("request time expired")
	}

	notCheckStr := fmt.Sprintf("%s%s%s", key, url, expTime)
	h := md5.New()
	h.Write([]byte(notCheckStr))
	cipherStr := h.Sum(nil)
	sign := hex.EncodeToString(cipherStr)
	//fmt.Println("req secret: ", secret, "sign", sign, "notCheckStr: ", notCheckStr)
	if secret != sign {
		return errors.New("auth failed")
	}
	return nil
}
