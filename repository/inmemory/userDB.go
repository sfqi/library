package inmemory

import (
	"fmt"
	"github.com/sfqi/library/domain/model"
	"time"
)

var users = []model.User{
	{
		Id:        1,
		Email:     "prvi@email.com",
		Name:      "ime",
		LastName:  "prezime",
		CreatedAt: time.Now().Add(-10 * time.Second),
		UpdatedAt: time.Now().Add(-8 * time.Second),
	},
	{
		Id:        2,
		Email:     "drugi@email.com",
		Name:      "ime2",
		LastName:  "prezime2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

type UserDb struct{
	Id int
	users []model.User
}

func NewUserDb()*UserDb{
	return &UserDb{
		Id:    len(users),
		users: users,
	}
}

func(udb *UserDb) FindUserById(id int)(*model.User,error){
	user,_,err := udb.findUserById(id);
	return user,err
}

func (udb *UserDb) findUserById(id int)(*model.User,int,error){
	for i, u := range udb.users {
		if u.Id == id {
			return &u, i, nil
		}
	}
	return nil, -1, fmt.Errorf("error while finding User")
}

func(udb *UserDb) FindAllUsers()([]*model.User,error){
	pointers := make([]*model.User,len(udb.users))
	for i:=0;i<len(udb.users);i++{
		pointers[i] = &udb.users[i]
	}
	fmt.Println(pointers)
	return pointers,nil
}


func (udb *UserDb) CreateUser(user *model.User) error {
	udb.Id++
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	user.Id = udb.Id
	udb.users = append(udb.users, *user)
	return nil
}


func (udb *UserDb) UpdateUser(toUpdate *model.User) error {
	user, index, err := udb.findUserById(toUpdate.Id)
	if err != nil {
		return err
	}
	user.UpdatedAt = time.Now()
	user.Name = toUpdate.Name
	user.LastName = toUpdate.LastName
	toUpdate = user
	udb.users[index] = *user
	return nil
}

func (udb *UserDb) DeleteUser(user *model.User) error {
	_, loc, err := udb.findUserById(user.Id)
	if err != nil {
		return err
	}
	udb.users = append(udb.users[:loc], udb.users[loc+1:]...)
	return nil
}
