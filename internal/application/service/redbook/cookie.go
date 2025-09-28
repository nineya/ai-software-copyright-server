package redbook

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"sync"
	"xorm.io/xorm"
)

type CookieService struct {
	service.AdminCrudService[table.RedbookCookie]
}

var onceCookie = sync.Once{}
var cookieService *CookieService

// 获取单例
func GetCookieService() *CookieService {
	onceCookie.Do(func() {
		cookieService = new(CookieService)
		cookieService.Db = global.DB
	})
	return cookieService
}

func (s *CookieService) CreateInBatch(adminId int64, param []table.RedbookCookie) error {
	return s.DbTransaction(func(session *xorm.Session) error {
		xhsUserIds := utils.ListTransform(param, func(item table.RedbookCookie) string {
			return item.XhsUserId
		})
		_, err := s.WhereAdminSession(adminId).In("xhs_user_id", xhsUserIds).Delete(&table.RedbookCookie{})
		if err != nil {
			return err
		}
		for i := range param {
			param[i].AdminId = adminId
		}
		_, err = s.Db.Insert(param)
		return err
	})
}

func (s *CookieService) UpdateStatusById(adminId int64, id int64, status enum.CookieStatus) error {
	_, err := s.WhereAndOmitAdminSession(adminId).ID(id).Update(&table.RedbookCookie{Status: status})
	return err
}

func (s *CookieService) InnerUpdateStatusById(id int64, status enum.CookieStatus) error {
	_, err := s.Db.ID(id).Update(&table.RedbookCookie{Status: status})
	return err
}

func (s *CookieService) GetAllByNormal() ([]table.RedbookCookie, error) {
	list := make([]table.RedbookCookie, 0)
	err := s.Db.Where("status = ?", enum.CookieStatus(1)).Find(&list)
	return list, err
}

func (s *CookieService) GetRand() ([]table.RedbookCookie, error) {
	list := make([]table.RedbookCookie, 0)
	err := s.Db.Limit(5, 0).OrderBy("RAND()").Find(&list)
	return list, err
}
