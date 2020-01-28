package main

import (
	"fmt"
	"github.com/sfqi/library/service"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sfqi/library/interactor"
	"github.com/sfqi/library/repository/postgres"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sfqi/library/log"
	"github.com/sfqi/library/openlibrary"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/handler"
	"github.com/sfqi/library/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	olc := openlibrary.NewClient(os.Getenv("OPEN_LIBRARY_URL"))
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
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
	bookInteractor := interactor.NewBook(store, olc)
	bookHandler := &handler.BookHandler{
		Interactor: bookInteractor,
	}

	uuidGenerator := &service.Generator{}
	loanInteractor := interactor.NewLoan(store)
	loanHandler := &handler.LoanHandler{
		Interactor: loanInteractor,
	}

	bookLoanInteractor := interactor.NewBookLoan(store, uuidGenerator)
	newLoanHandler := &handler.NewLoanHandler{Interactor: bookLoanInteractor}

	logger := log.New()

	bodyDump := middleware.BodyDump{
		Logger: logger,
	}

	handleFunc := handler.ErrorHandler{Logger: logger}.Wrap

	bookLoad := middleware.BookLoader{
		Interactor: bookInteractor,
	}
	r := mux.NewRouter()
	s := r.PathPrefix("/books").Subrouter()

	r.Handle("/books", handleFunc(bookHandler.Index)).Methods("GET")
	r.Handle("/books", handleFunc(bookHandler.Create)).Methods("POST")
	s.Handle("/{id}", handleFunc(bookHandler.Update)).Methods("PUT")
	s.Handle("/{id}", handleFunc(bookHandler.Get)).Methods("GET")
	s.Handle("/{id}", handleFunc(bookHandler.Delete)).Methods("DELETE")

	//loans endpoints
	s.Handle("/{id}/loans", handleFunc(loanHandler.FindLoansByBookID)).Methods("GET")
	s.Handle("/{id}/borrow", handleFunc(newLoanHandler.BorrowBook)).Methods("GET")

	s.Handle("/{id}/return", handleFunc(newLoanHandler.ReturnBook)).Methods("GET")

	r.Use(bodyDump.Dump)
	s.Use(bookLoad.GetBook)

	http.ListenAndServe(":8080", r)
}
