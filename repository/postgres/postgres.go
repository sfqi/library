package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/library/domain/model"
	"gopkg.in/yaml.v2"
	"os"
)


type dbStore struct {
	DB *gorm.DB
}

type Config struct{
	Server struct{
		Port int `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct{
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
		Dbname string 	`yaml:"dbname"`
	}`yaml:"database"`
}
func LoadConfig(name string)(*Config,error){
	f,err := os.Open(name)
	if err != nil{
		panic(err)
	}
	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil,err
	}
	return &cfg,nil
}


func Open(config Config)(*dbStore, error){
	store := dbStore{}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Server.Host,config.Server.Port,config.Database.Username,config.Database.Password,config.Database.Dbname)
	db, err := gorm.Open("postgres",psqlInfo)
	if err != nil {
		panic(err)
		return nil,err
	}
	store.DB = db
	return &store,nil
}
func(db *dbStore)FindById(id int)(*model.Book, error){
	b := model.Book{}
	if err := db.DB.First(&b,id).Error;err!=nil{
		return nil,err
	}
	return &b,nil
}
