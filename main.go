package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sfqi/library/repository/postgres"
	"net/http"
	"os"

	"github.com/sfqi/library/log"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sfqi/library/openlibrary"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/handler"
	"github.com/sfqi/library/middleware"
)

func init(){
	e := godotenv.Load()
	if e != nil{
		panic(e)
	}
	fmt.Println("Successfully read .env file")
}

func main() {
	openLibraryUrl := os.Getenv("LIBRARY")
	olc := openlibrary.NewClient(openLibraryUrl)
	config := postgres.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_DATABASE"),
	}

	//db := inmemory.NewDB()

	store ,err:= postgres.Open(config)
	if err != nil{
		panic(err)
	}
	fmt.Println("Successfully connected")
	bookHandler := &handler.BookHandler{
		Db: store,
		Olc: olc,
	}

	bodyDump := middleware.BodyDump{
		Logger: log.New(),
	}
	r := mux.NewRouter()

	r.HandleFunc("/books", bookHandler.Index).Methods("GET")
	r.HandleFunc("/books", bookHandler.Create).Methods("POST")
	r.HandleFunc("/books/{id}", bookHandler.Update).Methods("PUT")
	r.HandleFunc("/book/{id}", bookHandler.Get).Methods("GET")
	r.HandleFunc("/book/{id}", bookHandler.Delete).Methods("DELETE")
	r.Use(bodyDump.Dump)
	http.ListenAndServe(":8080", r)
}
