package interactor

import "github.com/sfqi/library/domain/model"

type userStore interface {
	FindUserByID(int) (*model.User, error)
}

type User struct {
	store userStore
}

func NewUser(store userStore) *User {
	return &User{
		store: store,
	}
}

func (u *User) FindByID(id int) (*model.User, error) {
	return u.store.FindUserByID(id)
}
