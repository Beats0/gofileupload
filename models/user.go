package models

import (
	"github.com/beats0/gofileupload/utils"
	"github.com/jinzhu/gorm"
)

/*
用户表
*/

type User struct {
	UId    uint   `gorm:"primary_key;uid;COMMENT:'uid';size:11;AUTO_INCREMENT"`
	Uname  string `gorm:"uname;COMMENT:'用户名';size:15;"`
	Avatar string `gorm:"avatar;COMMENT:'用户头像';size:50;"`
	Date   int64  `gorm:"date;COMMENT:'注册时间';size:10;"`
	Mail   string `gorm:"mail;COMMENT:'邮箱';size:32;"`
	Pwd    string `gorm:"pwd;COMMENT:'密码(加密)';size:32;"`
}

func SaveUser(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UserLogin(mail, pwd string) (*User, error) {
	var user User
	if err := db.
		First(&user, "mail=?", mail).
		Error; err != nil {
		return nil, err
	}
	if user.Pwd == utils.PwdSaltMd5(pwd) {
		return &user, nil
	}
	return nil, nil
}

// 邮箱是否被注册
func IsMailExist(mail string) (bool, error) {
	var total int
	err := db.Table("user").Where("mail=?", mail).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if total == 0 {
		return false, nil
	}
	return true, nil
}

// 用户是否存在
func IsUserExist(uid string) bool {
	var user User
	if err := db.Table("user").Where("uid=?", uid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return true
	}
	return false
}

func GetUserByMail(mail string) (*User, error) {
	var user User
	err := db.Table("user").Where("mail=?", mail).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &user, nil
}

func GetUserByUid(uid uint) (*User, error) {
	var user User
	if err := db.
		First(&user, "uid=?", uid).
		Select("uid uname avatar").
		Error; err != nil {
		return nil, err
	}
	return &user, nil
}
