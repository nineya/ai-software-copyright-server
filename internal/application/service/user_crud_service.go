package service

import (
	"ai-software-copyright-server/internal/application/model/table"
	"reflect"
)

type UserCrudService[T any] struct {
	BaseService
}

func (s *UserCrudService[T]) Create(userId int64, param T) (T, error) {
	reflect.ValueOf(&param).Interface().(table.UserTable).SetUserId(userId)
	_, err := s.Db.Insert(&param)
	return param, err
}

func (s *UserCrudService[T]) CreateInBatch(userId int64, param []T) error {
	for i := range param {
		reflect.ValueOf(&param[i]).Interface().(table.UserTable).SetUserId(userId)
	}
	_, err := s.Db.Insert(param)
	return err
}

func (s *UserCrudService[T]) DeleteById(userId int64, id int64) error {
	_, err := s.WhereUserSession(userId).ID(id).Delete(new(T))
	return err
}

func (s *UserCrudService[T]) DeleteInBatch(userId int64, ids []int64) error {
	_, err := s.WhereUserSession(userId).In("id", ids).Delete(new(T))
	return err
}

func (s *UserCrudService[T]) UpdateById(userId int64, id int64, param T) error {
	_, err := s.WhereAndOmitUserSession(userId).ID(id).AllCols().Update(&param)
	return err
}

func (s *UserCrudService[T]) GetById(userId, id int64) (*T, error) {
	mod := new(T)
	_, err := s.WhereUserSession(userId).ID(id).Get(mod)
	return mod, err
}

func (s *UserCrudService[T]) GetAll(userId int64) ([]T, error) {
	list := make([]T, 0)
	err := s.WhereUserSession(userId).Find(&list)
	return list, err
}
