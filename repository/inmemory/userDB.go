package inmemory

import (
	"fmt"
	"github.com/sfqi/library/domain/model"
	"time"
)

func (db *DB) FindUserById(id int) (*model.User, error) {
	user, _, err := db.findUserById(id)
	return user, err
}

func (db *DB) findUserById(id int) (*model.User, int, error) {
	for i, u := range db.users {
		if u.Id == id {
			return &u, i, nil
		}
	}
	return nil, -1, fmt.Errorf("error while finding User")
}

func (db *DB) FindAllUsers() ([]*model.User, error) {
	pointers := make([]*model.User, len(db.users))
	for i := 0; i < len(db.users); i++ {
		pointers[i] = &db.users[i]
	}
	fmt.Println(pointers)
	return pointers, nil
}

func (db *DB) CreateUser(user *model.User) error {
	db.Id++
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	user.Id = db.Id
	db.users = append(db.users, *user)
	return nil
}

func (db *DB) UpdateUser(toUpdate *model.User) error {
	user, index, err := db.findUserById(toUpdate.Id)
	if err != nil {
		return err
	}
	user.UpdatedAt = time.Now()
	user.Name = toUpdate.Name
	user.LastName = toUpdate.LastName
	toUpdate = user
	db.users[index] = *user
	return nil
}

func (db *DB) DeleteUser(user *model.User) error {
	_, loc, err := db.findUserById(user.Id)
	if err != nil {
		return err
	}
	db.users = append(db.users[:loc], db.users[loc+1:]...)
	return nil
}
