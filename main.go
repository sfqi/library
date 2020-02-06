package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/sfqi/library/service"

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
	userInteractor := interactor.NewUser(store)
	bookHandler := &handler.BookHandler{
		Interactor: bookInteractor,
	}

	uuidGenerator := &service.Generator{}
	loanInteractor := interactor.NewLoan(store)
	readLoanHandler := &handler.ReadLoanHandler{
		Interactor: loanInteractor,
	}

	bookLoanInteractor := interactor.NewBookLoan(store, uuidGenerator)
	writeLoanHandler := &handler.WriteLoanHandler{Interactor: bookLoanInteractor}

	logger := log.New()

	bodyDump := middleware.BodyDump{
		Logger: logger,
	}

	handleFunc := handler.ErrorHandler{Logger: logger}.Wrap

	bookLoad := middleware.BookLoader{
		Interactor: bookInteractor,
		Logger:     logger,
	}

	userLoad := middleware.UserLoader{
		Interactor: userInteractor,
	}

	r := mux.NewRouter()
	s := r.PathPrefix("/books").Subrouter()
	u := r.PathPrefix("/users").Subrouter()

	r.Handle("/books", handleFunc(bookHandler.Index)).Methods("GET")
	r.Handle("/books", handleFunc(bookHandler.Create)).Methods("POST")
	s.Handle("/{id}", handleFunc(bookHandler.Update)).Methods("PUT")
	s.Handle("/{id}", handleFunc(bookHandler.Get)).Methods("GET")
	s.Handle("/{id}", handleFunc(bookHandler.Delete)).Methods("DELETE")

	//loans endpoints
	s.Handle("/{id}/loans", handleFunc(readLoanHandler.FindLoansByBookID)).Methods("GET")
	s.Handle("/{id}/borrow", handleFunc(writeLoanHandler.BorrowBook)).Methods("POST")

	s.Handle("/{id}/return", handleFunc(writeLoanHandler.ReturnBook)).Methods("POST")

	u.Handle("/{id}/loans", handleFunc(readLoanHandler.FindLoansByUserID)).Methods("GET")

	r.Use(bodyDump.Dump)
	s.Use(bookLoad.GetBook)
	u.Use(userLoad.GetUser)

	http.ListenAndServe(":8080", r)
}
