package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/beats0/gofileupload/conf"
)

func PwdSaltMd5(str string) string {
	str = str + conf.AppSetting.PwdSalt
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
