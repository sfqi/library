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
	Db db
}
type db struct{
	DB *gorm.DB
}

func NewStore()*Store{
	return &Store{
		Db:db{
			DB:&gorm.DB{},
		},
	}
}

func Open(config PostgresConfig)(*Store, error){
	Store := NewStore()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,config.Port,config.User,config.Password,config.Name)
	db, err := gorm.Open("postgres",psqlInfo)
	if err != nil {
		return nil,err
	}
	Store.Db.DB = db
	return Store,nil
}
func(store *Store)FindById(id int)(*model.Book, error){
	b := model.Book{}
	if err := store.Db.DB.First(&b,id).Error;err!=nil{
		return nil,err
	}
	return &b,nil
}
