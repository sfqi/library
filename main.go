package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sfqi/library/repository/postgres"
	"net/http"
	"os"
	"strconv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sfqi/library/log"
	"github.com/sfqi/library/openlibrary"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/handler"
	"github.com/sfqi/library/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil{
		panic(err)
	}

	olc := openlibrary.NewClient(os.Getenv("OPEN_LIBRARY_URL"))
	port,err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil{
		panic(err)
	}
	config := postgres.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	store, err := postgres.Open(config)
	if err != nil {
		panic(err)
	}
	defer store.Close()
	fmt.Println("Successfully connected")
	bookHandler := &handler.BookHandler{
		Db:  store,
		Olc: olc,
	}

	bodyDump := middleware.BodyDump{
		Logger: log.New(),
	}

	bookLoad := middleware.BookLoader{
		Db: store,
	}
	r := mux.NewRouter()
	s := r.PathPrefix("/books").Subrouter()

	r.HandleFunc("/books", bookHandler.Index).Methods("GET")
	r.HandleFunc("/books", bookHandler.Create).Methods("POST")
	s.HandleFunc("/{id}", bookHandler.Update).Methods("PUT")
	s.HandleFunc("/{id}", bookHandler.Get).Methods("GET")
	s.HandleFunc("/{id}", bookHandler.Delete).Methods("DELETE")
	r.Use(bodyDump.Dump)
	s.Use(bookLoad.GetBook)

	http.ListenAndServe(":8080", r)
}
