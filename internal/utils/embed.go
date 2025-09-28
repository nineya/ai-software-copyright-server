package utils

import (
	"ai-software-copyright-server/internal/global"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func CopyEmbedDir(sourcePath, targetPath string) error {
	files, err := global.FS.ReadDir(sourcePath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(targetPath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "创建文件目录失败")
	}
	for _, file := range files {
		s := sourcePath + "/" + file.Name()
		t := targetPath + "/" + file.Name()
		if file.IsDir() {
			err = CopyEmbedDir(s, t)
			if err != nil {
				return err
			}
		} else {
			err = CopyEmbedFile(s, t)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyEmbedFile(sourceFile, targetFile string) error {
	fileContent, err := global.FS.ReadFile(sourceFile)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(targetFile, fileContent, 0644)
}

// 读取文件内容
func ReadEmbedFileContent(sourceFile string) (string, error) {
	fileContent, err := global.FS.ReadFile(sourceFile)
	if err != nil {
		return "", err
	}
	return string(fileContent), err
}
