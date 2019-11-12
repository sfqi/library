package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/library/domain/model"
)

type PostgresConfig struct{
	Host 		string
	Port 		int
	User 		string
	Password 	string
	Name 		string
}


type Store struct {
	db *gorm.DB
}

func NewStore()*Store{
	return &Store{}
}

func Open(config PostgresConfig)(*Store, error){
	Store := NewStore()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,config.Port,config.User,config.Password,config.Name)

	db, err := gorm.Open("postgres",psqlInfo)
	if err != nil {
		return nil,err
	}

	Store.db = db
	return Store,nil
}
func(store *Store)FindById(id int)(*model.Book, error){
	b := model.Book{}
	if err := store.db.First(&b,id).Error;err!=nil{
		return nil,err
	}
	return &b,nil
}
