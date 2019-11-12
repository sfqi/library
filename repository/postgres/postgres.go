package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sfqi/library/domain/model"
)
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type dbStore struct {
	DB *gorm.DB
}

func Open(config PostgresConfig) (*dbStore, error) {
	store := dbStore{}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
		return nil, err
	}
	store.DB = db
	return &store, nil
}
func (db *dbStore) FindById(id int) (*model.Book, error) {
	b := model.Book{}
	if err := db.DB.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}
