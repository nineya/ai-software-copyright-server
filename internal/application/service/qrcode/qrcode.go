package qrcode

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	attaSev "ai-software-copyright-server/internal/application/service/attachment"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"mime/multipart"
	"strings"
	"sync"
)

type QrcodeService struct {
	service.UserCrudService[table.Qrcode]
}

var onceQrcode = sync.Once{}
var qrcodeService *QrcodeService

// 获取单例
func GetQrcodeService() *QrcodeService {
	onceQrcode.Do(func() {
		qrcodeService = new(QrcodeService)
		qrcodeService.Db = global.DB
	})
	return qrcodeService
}

func (s *QrcodeService) Build(userId int64, param request.QrcodeBuildParam) (*response.QrcodeBuildResponse, error) {
	expenseCredits := 10
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	code, err := qrcode.Encode(param.Content, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	result := &response.QrcodeBuildResponse{Content: base64.StdEncoding.EncodeToString(code)}

	// 扣款
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(9), expenseCredits, fmt.Sprintf("购买二维码生成服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits

	return result, nil
}

// 添加图片
func (s *QrcodeService) AddImageById(userId, id int64, file *multipart.FileHeader) (*response.QrcodeAddImageResponse, error) {
	var bytes []byte
	// 预检查图片
	if file.Size > 200*1024 {
		// 图片压缩
		f, err := file.Open() // 读取文件
		if err != nil {
			return nil, errors.Wrap(err, "读取图片失败")
		}
		defer f.Close()
		bytes, err = utils.ImageCompress(f)
		if len(bytes) > 200*1024 {
			return nil, errors.New("上传图片太大了，请压缩")
		}
	}

	expenseCredits := 30
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	mod, err := s.GetById(userId, id)
	if err != nil {
		return nil, err
	}
	if mod.Id == 0 {
		return nil, errors.New("活码不存在")
	}
	// 保存图片
	var image *response.UploadImageResponse
	if bytes != nil && len(bytes) > 0 {
		image, err = attaSev.GetImageService().UploadByBytes(bytes, ".png", "qrcode")
	} else {
		image, err = attaSev.GetImageService().Upload(file, "qrcode")
	}
	if err != nil {
		return nil, err
	}
	if mod.TargetUrls == nil {
		mod.TargetUrls = []string{image.Url}
	} else {
		mod.TargetUrls = append(mod.TargetUrls, image.Url)
	}
	err = s.UpdateById(userId, id, *mod)
	if err != nil {
		return nil, err
	}

	result := &response.QrcodeAddImageResponse{Qrcode: *mod}

	// 扣款
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(18), expenseCredits, fmt.Sprintf("购买活码添加图片服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits

	return result, nil
}

func (s *QrcodeService) DeleteById(userId int64, id int64) error {
	mod, err := s.GetById(userId, id)
	if err != nil {
		return err
	}
	if mod.Id == 0 {
		return nil
	}
	for _, url := range mod.TargetUrls {
		_ = attaSev.GetImageService().Delete(url)
	}
	_, err = s.WhereUserSession(userId).ID(id).Delete(mod)
	return err
}

// 删除图片
func (s *QrcodeService) DeleteImageById(userId, id int64, imageUrl string) error {
	imageUrl = strings.TrimSpace(imageUrl)
	if imageUrl == "" {
		return errors.New("图片不存在")
	}
	mod, err := s.GetById(userId, id)
	if err != nil {
		return err
	}
	if mod.Id == 0 {
		return errors.New("活码不存在")
	}

	// 删除文件
	mod.TargetUrls = utils.ListFilter(mod.TargetUrls, func(item string) bool {
		if imageUrl == item {
			_ = attaSev.GetImageService().Delete(imageUrl)
			return false
		}
		return true
	})

	_, err = s.WhereUserSession(userId).ID(id).AllCols().Update(mod)
	return err
}

// 分页查询
func (s *QrcodeService) GetByPage(userId int64, param request.PageableParam) (*response.PageResponse, error) {
	session := s.WhereUserSession(userId).Desc("create_time")
	list := make([]table.Qrcode, 0)
	return s.HandlePageable(param, &list, session)
}

func (s *QrcodeService) GetByAlias(alias string) (*table.Qrcode, error) {
	mod := &table.Qrcode{Alias: alias}
	_, err := s.Db.Get(mod)
	return mod, err
}

func (s *QrcodeService) UpdateVisitsIncreaseById(id int64, userAgent string) error {
	// 爬虫请求不处理，直接返回
	if global.BotReg.MatchString(userAgent) {
		return nil
	}
	_, err := s.Db.ID(id).Incr("visits", 1).NoAutoTime().Update(&table.Qrcode{})
	return err
}
