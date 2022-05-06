package gorm_expand

type Service[T any] struct {
	m Mapper[T]
}

func (s Service[T]) SetMapper(m Mapper[T]) {
	s.m = m
}

func (s Service[T]) Insert(t T) error {
	return s.m.Insert(t)
}

func (s Service[T]) UpdateById(t T) error {
	return s.m.UpdateById(t)
}

func (s Service[T]) DeleteById(id any) error {
	return s.m.DeleteById(id)
}

func (s Service[T]) GetById(id any) (T, error) {
	return s.m.GetById(id)
}

func (s Service[T]) Get(t T, condition queryCondition) (T, error) {
	return s.m.Get(t, condition)
}

func (s Service[T]) List() ([]T, error) {
	return s.m.List()
}

func (s Service[T]) ListByEntity(t T, condition queryCondition) ([]T, error) {
	return s.m.ListByCondition(t, condition)
}
