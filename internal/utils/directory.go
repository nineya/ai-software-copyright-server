package utils

import (
	"ai-software-copyright-server/internal/global"
	"fmt"
	"os"
	"path/filepath"
)

// @description: 文件目录是否存在
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

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// 取得图片文件物理存储路径
func GetImageStorePath() string {
	return fmt.Sprintf("%s/data/image", global.WORK_DIR)
}

// 取得图片文件物理存储路径
func GetSoftwareCopyrightPath(id int64) string {
	return fmt.Sprintf("%s/data/software_copyright/%d", global.WORK_DIR, id)
}
