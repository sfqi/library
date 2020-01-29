package model

import "time"

type User struct {
	Id        int
	Email     string
	Name      string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
