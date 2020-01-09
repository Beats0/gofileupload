package utils

import (
	"encoding/binary"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// 文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 检查文件夹并创建
func CheckDirPath(path string) (bool, error) {
	exist, _ := PathExists(path)
	if exist {
		return true, nil
	} else {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

// 获取part文件
func GetFilePart(dirPth string) (files []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		ok := strings.HasSuffix(fi.Name(), ".part")
		if ok {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	// 注意将part文件按文件名排序, 依次合并
	sort.Slice(
		files,
		func(i, j int) bool {
			return SortName(files[i]) < SortName(files[j])
		},
	)
	return files, nil
}

// sortName returns a filename sort key with
// non-negative integer suffixes in numeric order.
// For example, amt, amt0, amt2, amt10, amt099, amt100, ...
func SortName(filename string) string {
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	// split numeric suffix
	i := len(name) - 1
	for ; i >= 0; i-- {
		if '0' > name[i] || name[i] > '9' {
			break
		}
	}
	i++
	// string numeric suffix to uint64 bytes
	// empty string is zero, so integers are plus one
	b64 := make([]byte, 64/8)
	s64 := name[i:]
	if len(s64) > 0 {
		u64, err := strconv.ParseUint(s64, 10, 64)
		if err == nil {
			binary.BigEndian.PutUint64(b64, u64+1)
		}
	}
	// prefix + numeric-suffix + ext
	return name[:i] + string(b64) + ext
}

// 文件是否存在
func FileExit(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 获取文件夹大小
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
