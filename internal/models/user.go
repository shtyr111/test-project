package models

import (
	"github.com/google/uuid"
)

type User struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Age    int       `json:"age"`
	Number int       `json:"number"`
	Status string    `json:"status"`
}

func (user *User) SetStatus(Status string) {
	user.Status = Status
}
