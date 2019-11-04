package postgres


import(
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"


)
const(
	host ="localhost"
	port =5432
	user="bojan"
	password="bojan"
	dbname="library"
)
//this struct is just used for demonstration
type Book struct{
	gorm.Model
	Title string
	Year string
}

type Store interface {
	FindById(id int)(*Book, error)
}

type dbStore struct {
	Db *gorm.DB
}

func Open()(*dbStore, error){
	store := dbStore{}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,port,user,password,dbname)
	db, err := gorm.Open("postgres",psqlInfo)
	if err != nil {
		panic(err)
		return nil,err
	}
	store.Db = db
	return &store,nil

}
func(db *dbStore)FindById(id int)(*Book, error){
	b := Book{}
	if err := db.Db.Where("id=?",id).First(&b).Error;err!=nil{
		return nil,err
	}
	return &b,nil
}
