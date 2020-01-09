package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

/*
上传文件表
*/

type Upload struct {
	ID            uint   `gorm:"primary_key;COMMENT:'id';size:11;AUTO_INCREMENT;" json:"id"`
	UId           uint   `gorm:"uid;COMMENT:'上传uid';size:11;" json:"uid"`
	Date          int64  `gorm:"date;COMMENT:'上传时间';size:10;" json:"date"`
	File_name     string `gorm:"file_name;COMMENT:'文件名称';size:255;" json:"file_name"`
	File_size     int64  `gorm:"file_size;COMMENT:'文件大小';size:11;" json:"file_size"`
	File_ext      string `gorm:"file_ext;COMMENT:'文件后缀';size:11;" json:"file_ext"`
	File_type     string `gorm:"file_type;COMMENT:'文件类型';size:11;" json:"file_type"`
	Is_dir        int    `gorm:"is_dir;COMMENT:'是否为dir';DEFAULT: 0" json:"is_dir"`
	Parent        uint   `gorm:"parent;COMMENT:'父id'" json:"parent"`
	Md5           string `gorm:"md5;COMMENT:'md5';size:32;" json:"md5"`
	Last_modified int64  `gorm:"last_modified;COMMENT:'上次更改时间';size:10;" json:"last_modified"`
}

type UserFile struct {
	ID            uint   `json:"id"`
	UId           uint   `json:"uid"`
	Date          int64  `json:"date"`
	File_name     string `json:"file_name"`
	File_size     int64  `json:"file_size"`
	File_ext      string `json:"file_ext"`
	File_type     string `json:"file_type"`
	Md5           string `json:"md5"`
	Is_dir        int    `json:"is_dir"`
	Last_modified string `json:"last_modified"`
}

func SaveUploadFile(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

// 文件列表
func GetUserUploadByUid(uid uint, dirId uint, orderBy string, desc string, page, limit int) ([]*UserFile, error) {
	var userFiles []*UserFile
	orderQuery := fmt.Sprintf("%s %s", orderBy, desc)
	err := db.Table("upload").
		Where("uid = ? AND parent= ?", uid, dirId).
		Order(orderQuery).
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&userFiles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return userFiles, nil
}

// 对应项目文件total
func GetUserUploadTotalByUid(dirId uint, uid uint) (int, error) {
	var total int
	err := db.Table("upload").
		Where("parent=? AND uid=?", dirId, uid).
		Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return total, nil
}

// 文件分类列表
func GetUserFileTypesByUid(uid uint, fileType string, orderBy string, desc string, page, limit int) ([]*UserFile, error) {
	var userFiles []*UserFile
	orderQuery := fmt.Sprintf("%s %s", orderBy, desc)
	err := db.Table("upload").
		Where("uid=? AND file_type=?", uid, fileType).
		Order(orderQuery).
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&userFiles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return userFiles, nil
}

// 对应分类total
func GetUserFileTypeTotal(fileType string, uid uint) (int, error) {
	var total int
	err := db.Table("upload").
		Where("uid=? AND file_type=?", uid, fileType).
		Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return total, nil
}

// 当前文件夹下该文件
func GetUserFileByDirAndFileName(uid uint, dirId uint, fileName string) ([]*UserFile, error) {
	var userFiles []*UserFile
	err := db.Table("upload").
		Where("parent=? AND uid=? AND file_name=?", dirId, uid, fileName).
		Scan(&userFiles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return userFiles, nil
}

// 检查文件所属
func CheckUserFile(fsId uint, uid uint) (UserFile, error) {
	var userFile UserFile
	if err := db.Table("upload").
		Where("id=? AND uid=?", fsId, uid).
		First(&userFile).Error; err != nil {
		return userFile, err
	}
	return userFile, nil
}

// 检查文件所属
func CheckUserFileMd5(md5 string, uid uint) (UserFile, error) {
	var userFile UserFile
	if err := db.Table("upload").
		Where("md5=? AND uid=?", md5, uid).
		First(&userFile).Error; err != nil {
		return userFile, err
	}
	return userFile, nil
}

// 删除文件记录
func DeleteUserFile(fsId uint, uid uint) error {
	if err := db.Table("upload").
		Where("id=? AND uid=?", fsId, uid).
		Delete(UserFile{}).Error; err != nil {
		return err
	}
	return nil
}

// 检查文件MD5
func CheckFile(md5 string, uid uint) (UserFile, error) {
	var userFile UserFile
	if err := db.Table("upload").
		Where("md5=? AND uid=?", md5, uid).
		First(&userFile).Error; err != nil {
		return userFile, err
	}
	return userFile, nil
}

// 创建文件夹
func CreateFolder(dirData Upload) (Upload, error) {
	if err := db.Save(&dirData).Error; err != nil {
		return dirData, err
	}
	return dirData, nil
}

// 检查同一文件夹下是否含有该文件
func IsItemExist(dirId uint, uid uint, fileName string) (int, error) {
	var total int
	err := db.Table("upload").
		Where("parent=? AND uid=? AND file_name=?", dirId, uid, fileName).
		Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return total, nil
}

// 查看文件所在位置
func FindFilePath(fsId uint, uid uint) (uint, error) {
	var upload Upload
	if err := db.Table("upload").
		Where("id=? AND uid=?", fsId, uid).
		First(&upload).Error; err != nil {
		return upload.Parent, err
	}
	return upload.Parent, nil
}

// 文件搜索
func Search(uid uint, q string, orderBy string, desc string, page, limit int) ([]*UserFile, error) {
	var userFiles []*UserFile
	orderQuery := fmt.Sprintf("%s %s", orderBy, desc)
	err := db.Table("upload").
		Where("uid=? AND file_name LIKE ?", uid, "%"+q+"%").
		Order(orderQuery).
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&userFiles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return userFiles, nil
}

// 重命名
func Rename(fsId uint, uid uint, fileName string) error {
	if err := db.Table("upload").
		Where("id=? AND uid=?", fsId, uid).
		Update("file_name", fileName).Error; err != nil {
		return err
	}
	return nil
}

// 文件搜索 total
func SearchTotal(uid uint, q string) (int, error) {
	var total int
	err := db.Table("upload").
		Where("uid=? AND file_name LIKE ?", uid, "%"+q+"%").
		Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return total, nil
}

func GetFileByMd5(md5 string) (UserFile, error) {
	var userFile UserFile
	if err := db.Table("upload").
		Where("md5=?", md5).
		First(&userFile).Error; err != nil {
		return userFile, err
	}
	return userFile, nil
}

// 视频信息
func GetVideoInfo(uid uint, md5 string) (UserFile, error) {
	var userFile UserFile
	if err := db.Table("upload").
		Where("uid=? AND md5=?", uid, md5).
		First(&userFile).Error; err != nil {
		return userFile, err
	}
	return userFile, nil
}
