package service

import (
	"ai-software-copyright-server/internal/application/model/table"
	"reflect"
)

type AdminCrudService[T any] struct {
	BaseService
}

func (s *AdminCrudService[T]) Create(adminId int64, param T) (T, error) {
	reflect.ValueOf(&param).Interface().(table.AdminTable).SetAdminId(adminId)
	_, err := s.Db.Insert(&param)
	return param, err
}

func (s *AdminCrudService[T]) CreateInBatch(adminId int64, param []T) error {
	for i := range param {
		reflect.ValueOf(&param[i]).Interface().(table.AdminTable).SetAdminId(adminId)
	}
	_, err := s.Db.Insert(param)
	return err
}

func (s *AdminCrudService[T]) DeleteById(adminId int64, id int64) error {
	_, err := s.WhereAdminSession(adminId).ID(id).Delete(new(T))
	return err
}

func (s *AdminCrudService[T]) DeleteInBatch(adminId int64, ids []int64) error {
	_, err := s.WhereAdminSession(adminId).In("id", ids).Delete(new(T))
	return err
}

func (s *AdminCrudService[T]) UpdateById(adminId int64, id int64, param T) error {
	_, err := s.WhereAndOmitAdminSession(adminId).ID(id).AllCols().Update(&param)
	return err
}

func (s *AdminCrudService[T]) GetById(adminId, id int64) (*T, error) {
	mod := new(T)
	_, err := s.WhereAdminSession(adminId).ID(id).Get(mod)
	return mod, err
}

func (s *AdminCrudService[T]) GetAll(adminId int64) ([]T, error) {
	list := make([]T, 0)
	err := s.WhereAdminSession(adminId).Find(&list)
	return list, err
}
