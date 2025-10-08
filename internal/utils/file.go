package utils

import (
	"archive/zip"
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func GetFileMimeType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", errors.Wrap(err, "打开文件失败")
	}
	defer file.Close()

	// 读取文件的前 512 个字节用于检测 MIME 类型
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", errors.Wrap(err, "读取文件失败")
	}

	return http.DetectContentType(buffer), nil
}

func ImageCompress(file multipart.File) ([]byte, error) {
	image, err := imaging.Decode(file)
	if err != nil {
		return nil, errors.Wrap(err, "读取图片失败")
	}
	bounds := image.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	//生成缩略图，压缩80%
	image = imaging.Resize(image, int(float64(width)*0.8), int(float64(height)*0.8), imaging.Linear)

	// 保存图片
	var buf bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	err = encoder.Encode(&buf, image)
	if err != nil {
		return nil, errors.Wrap(err, "图片压缩失败")
	}

	return buf.Bytes(), err
}

// 使用 image/png 解码（推荐）
func ImagePngWithDecode(data []byte) (width, height int, err error) {
	// 创建 bytes.Reader
	reader := bytes.NewReader(data)

	// 解码 PNG 配置信息（不解码像素数据，性能更好）
	config, err := png.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}

	return config.Width, config.Height, nil
}

func CopyFile(file *multipart.FileHeader, dstDir, dstFile string) error {
	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "创建文件目录失败")
	}

	f, err := file.Open() // 读取文件
	if err != nil {
		return errors.Wrap(err, "打开文件失败")
	}
	defer f.Close()

	out, err := os.Create(filepath.Join(dstDir, dstFile))
	if err != nil {
		return errors.Wrap(err, "创建文件失败")
	}
	defer out.Close() // 创建文件 defer 关闭

	_, err = io.Copy(out, f)
	if err != nil {
		return errors.Wrap(err, "写入文件失败")
	}
	return nil
}

func WriteFile(bytes []byte, dstDir, dstFile string) error {
	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "创建文件目录失败")
	}

	out, err := os.Create(filepath.Join(dstDir, dstFile))
	if err != nil {
		return errors.Wrap(err, "创建文件失败")
	}
	defer out.Close() // 创建文件 defer 关闭

	_, err = out.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "写入文件失败")
	}
	return nil
}

func CreateZip(path, zipFile string) error {
	// 创建zip文件
	zipFileWriter, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zipFileWriter.Close()

	// 创建zip writer
	zipWriter := zip.NewWriter(zipFileWriter)
	defer zipWriter.Close()

	// 获取文件绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	return filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 创建zip文件中的文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 设置文件头中的文件名（相对路径）
		relPath, err := filepath.Rel(absPath, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// 如果是目录，需要在文件名后加斜杠
		if info.IsDir() {
			header.Name += "/"
		} else {
			// 设置压缩方法
			header.Method = zip.Deflate
		}

		// 创建zip文件中的文件
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果是文件，写入文件内容
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func UnZip(zipFile, dstPath string) error {
	archive, err := zip.OpenReader(zipFile)
	if err != nil {
		return errors.Wrap(err, "打开压缩文件失败")
	}
	defer archive.Close()

	for _, item := range archive.File {
		filePath := filepath.Join(dstPath, item.Name)
		if item.FileInfo().IsDir() {
			if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}
		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, item.Mode())
		defer file.Close()
		if err != nil {
			return err
		}
		fileInArchive, err := item.Open()
		defer fileInArchive.Close()
		if err != nil {
			return err
		}
		if _, err := io.Copy(file, fileInArchive); err != nil {
			return err
		}
	}
	return nil
}

func Move(srcPath, dstPath string) error {
	// 移动主题到对应目录，跨磁盘移动时这里会失败
	if err := os.Rename(srcPath, dstPath); err == nil {
		return nil
	}
	inputFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return err
	}
	os.Remove(srcPath)
	return nil

}

// 从后往前读取文件指定行
func ReadFileFromBack(filePath string, lineNum int) (string, error) {
	//打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	//获取文件大小
	fs, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fs.Size()

	var offset int64 = -1   //偏移量，初始化为-1，若为0则会读到EOF
	char := make([]byte, 1) //用于读取单个字节
	lineStr := ""           //存放一行的数据
	for (-offset) <= fileSize {
		//通过Seek函数从末尾移动游标然后每次读取一个字节
		file.Seek(offset, io.SeekEnd)
		_, err = file.Read(char)
		if err != nil {
			return "", err
		}
		if char[0] == '\n' {
			lineNum-- //到此读取完一行
			if lineNum == 0 {
				return lineStr, nil
			}
		}
		lineStr = string(char) + lineStr
		offset--
	}
	return lineStr, nil
}

// 将内容写入文件
func WriteToFile(path string, content []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}
	return nil
}
