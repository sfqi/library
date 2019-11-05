package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/library/domain/model"
	"os"
)
const(
	host ="localhost"
	port =5432
	dbname="library"
)

var username = os.Getenv("PSQLUSERNAME")
var password = os.Getenv("PSQLPASSWORD")

type dbStore struct {
	DB *gorm.DB
}

func Open()(*dbStore, error){
	store := dbStore{}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,port,username,password,dbname)
	db, err := gorm.Open("postgres",psqlInfo)
	if err != nil {
		panic(err)
		return nil,err
	}
	b := model.Book{}
	db.AutoMigrate(&b)
	store.DB = db
	return &store,nil

}
func(db *dbStore)FindById(id int)(*model.Book, error){
	b := model.Book{}
	if err := db.DB.Where("id=?",id).First(&b).Error;err!=nil{
		return nil,err
	}
	return &b,nil
}
