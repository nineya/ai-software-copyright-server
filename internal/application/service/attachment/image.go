package attachment

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/config"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ImageService struct {
	CONFIG config.ToolImage
}

var onceImage = sync.Once{}
var imageService *ImageService

// 获取单例
func GetImageService() *ImageService {
	onceImage.Do(func() {
		imageService = new(ImageService)
		imageService.CONFIG = global.CONFIG.Tool.ToolImage
	})
	return imageService
}

func (s *ImageService) Upload(file *multipart.FileHeader, bucket string) (*response.UploadImageResponse, error) {
	if bucket == "" {
		if !s.CONFIG.AllowUnknownBucket || s.CONFIG.DefaultBucket == "" {
			return nil, errors.New("未指定bucket名称")
		}
		bucket = s.CONFIG.DefaultBucket
	} else if !utils.ListContains(s.CONFIG.Buckets, bucket) {
		return nil, errors.New("不被允许的bucket名称：" + bucket)
	}

	// 上传的后缀名
	suffix := path.Ext(file.Filename)
	// 上传的文件名
	// name := strings.TrimSuffix(file.Filename, suffix)
	if suffix != ".png" && suffix != ".jpg" && suffix != ".jpeg" && suffix != ".gif" {
		return nil, errors.New("文件格式不符合要求: " + suffix)
	}
	// 生成存储的文件名和路径
	storeName := generateName() + suffix
	genPath := generatePath()

	err := utils.CopyFile(file, filepath.Join(utils.GetImageStorePath(), bucket, genPath), storeName)
	if err != nil {
		return nil, err
	}

	return &response.UploadImageResponse{
		Name: storeName,
		Url:  fmt.Sprintf("%s/%s%s%s", s.CONFIG.DownloadPath, bucket, genPath, storeName),
	}, err
}

func (s *ImageService) UploadByBytes(bytes []byte, suffix string, bucket string) (*response.UploadImageResponse, error) {
	if bucket == "" {
		if !s.CONFIG.AllowUnknownBucket || s.CONFIG.DefaultBucket == "" {
			return nil, errors.New("未指定bucket名称")
		}
		bucket = s.CONFIG.DefaultBucket
	} else if !utils.ListContains(s.CONFIG.Buckets, bucket) {
		return nil, errors.New("不被允许的bucket名称：" + bucket)
	}

	// 上传的文件名
	// name := strings.TrimSuffix(file.Filename, suffix)
	if suffix != ".png" && suffix != ".jpg" && suffix != ".jpeg" && suffix != ".gif" {
		return nil, errors.New("文件格式不符合要求: " + suffix)
	}
	// 生成存储的文件名和路径
	storeName := generateName() + suffix
	genPath := generatePath()

	err := utils.WriteFile(bytes, filepath.Join(utils.GetImageStorePath(), bucket, genPath), storeName)
	if err != nil {
		return nil, err
	}

	return &response.UploadImageResponse{
		Name: storeName,
		Url:  fmt.Sprintf("%s/%s%s%s", s.CONFIG.DownloadPath, bucket, genPath, storeName),
	}, err
}

func (s *ImageService) Delete(filepath string) error {
	strings.LastIndex(filepath, s.CONFIG.DownloadPath)
	index := strings.Index(filepath, "/image/")
	if index != -1 {
		filepath = utils.GetImageStorePath() + filepath[index+6:]
	}
	return os.Remove(filepath)
}

func generatePath() string {
	return time.Now().Format("/2006/01/")
}

func generateName() string {
	return time.Now().Format("02150405") + strconv.Itoa(int(uuid.New().ID()))
}
