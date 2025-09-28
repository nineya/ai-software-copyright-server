package service

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/global"
	"math"
	"xorm.io/xorm"
)

type BaseService struct {
	Db *xorm.Engine
}

// 处理事务的方法
func (s *BaseService) DbTransaction(exec func(*xorm.Session) error) (err error) {
	// 创建事务
	session := s.Db.NewSession()

	// 事务开始
	err = session.Begin()
	if err != nil {
		session.Close()
		return err
	}

	// 关闭事务
	defer func() {
		defer session.Close()
		if err != nil {
			global.LOG.Error("事务回滚：" + err.Error())
			// 事务回滚
			session.Rollback()
		} else if err2 := recover(); err2 != nil {
			// 事务回滚
			session.Rollback()
			panic(err2)
		} else {
			// 事务提交
			session.Commit()
		}
	}()

	// 事务相关操作
	return exec(session)
}

// 处理分页查询，page从0开始
func (s *BaseService) HandlePageable(pageable request.PageableParam, data interface{}, session *xorm.Session) (*response.PageResponse, error) {
	size := pageable.Size
	page := pageable.Page
	session.Limit(size, size*page)
	count, err := session.FindAndCount(data)
	if err != nil {
		return nil, err
	}
	return s.HandleContentPageable(data, count, pageable), nil
}

// 处理列表分页
func (s *BaseService) HandleContentPageable(content any, count int64, pageable request.PageableParam) *response.PageResponse {
	pages := int(math.Ceil(float64(count) / float64(pageable.Size)))
	return &response.PageResponse{
		Content:     content,
		HasPrevious: pageable.Page > 0,
		HasNext:     pageable.Page+1 < pages,
		Page:        pageable.Page,
		Pages:       pages,
		Total:       count,
	}
}

// 取得带 adminId 的数据库连接
func (s *BaseService) WhereAdminSession(adminId int64) *xorm.Session {
	return s.Db.Where("admin_id = ?", adminId)
}

func (s *BaseService) WhereAndOmitAdminSession(adminId int64) *xorm.Session {
	return s.Db.Where("admin_id = ?", adminId).Omit("admin_id")
}

func (s *BaseService) WhereAdminSessionByTable(tableName string, adminId int64) *xorm.Session {
	return s.Db.Table(tableName).Where(tableName+".admin_id = ?", adminId)
}

func (S *BaseService) AddWhereAdmin(adminId int64, session *xorm.Session) *xorm.Session {
	return session.Where("admin_id = ?", adminId)
}

func (S *BaseService) AddWhereAndOmitAdmin(adminId int64, session *xorm.Session) *xorm.Session {
	return session.Where("admin_id = ?", adminId).Omit("admin_id")
}

func (s *BaseService) AddWhereAdminByTable(tableName string, adminId int64, session *xorm.Session) *xorm.Session {
	return session.Table(tableName).Where(tableName+".admin_id = ?", adminId)
}

// 取得带 userId 的数据库连接
func (s *BaseService) WhereUserSession(userId int64) *xorm.Session {
	return s.Db.Where("user_id = ?", userId)
}

func (s *BaseService) WhereAndOmitUserSession(userId int64) *xorm.Session {
	return s.Db.Where("user_id = ?", userId).Omit("user_id")
}

func (s *BaseService) WhereUserSessionByTable(tableName string, userId int64) *xorm.Session {
	return s.Db.Table(tableName).Where(tableName+".user_id = ?", userId)
}

func (S *BaseService) AddWhereUser(userId int64, session *xorm.Session) *xorm.Session {
	return session.Where("user_id = ?", userId)
}

func (S *BaseService) AddWhereAndOmitUser(userId int64, session *xorm.Session) *xorm.Session {
	return session.Where("user_id = ?", userId).Omit("user_id")
}

func (s *BaseService) AddWhereUserByTable(tableName string, userId int64, session *xorm.Session) *xorm.Session {
	return session.Table(tableName).Where(tableName+".user_id = ?", userId)
}
