package models

import "github.com/jinzhu/gorm"

/*
上传表
*/

type Upload struct {
	ID        uint   `gorm:"primary_key;COMMENT:'id';size:11;AUTO_INCREMENT"`
	UId       uint   `gorm:"uid;COMMENT:'上传uid';size:11;"`
	Date      int64  `gorm:"date;COMMENT:'上传时间';size:11;"`
	File_name string `gorm:"file_name;COMMENT:'文件名称';size:255;"`
	File_size int64  `gorm:"file_size;COMMENT:'文件大小';size:11;"`
	File_ext  string `gorm:"file_ext;COMMENT:'文件后缀';size:11;"`
	File_type string `gorm:"file_type;COMMENT:'文件类型';size:11;"`
	Md5       string `gorm:"md5;COMMENT:'hash';size:32;"`
}

type UserFile struct {
	ID        uint   `json:"id"`
	UId       uint   `json:"uid"`
	Date      int64  `json:"date"`
	File_name string `json:"file_name"`
	File_size int64  `json:"file_size"`
	File_ext  string `json:"file_ext"`
	Md5       string `json:"md5"`
}

func SaveUploadFile(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func GetUserUploadByUid(uid uint, page, limit int) ([]*UserFile, error) {
	var userFiles []*UserFile
	sql := `SELECT * FROM upload WHERE upload.uid = ? ORDER BY upload.date ASC LIMIT ?, ?`
	err := db.Raw(sql, uid, (page-1)*limit, page*limit).Scan(&userFiles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return userFiles, nil
}

func GetUserUploadTotalByUid(uid uint) (int, error) {
	var total int
	err := db.Table("upload").Where("uid = ?", uid).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return total, nil
}
