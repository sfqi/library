package model

import "time"

type User struct {
	Id        int
	Email     string
	Name      string
	Lastname  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
