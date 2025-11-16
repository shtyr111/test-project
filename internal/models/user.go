package models

import (
	"github.com/google/uuid"
)

type User struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Age    int       `json:"age"`
	Number int       `json:"number"`
}

func (user *User) SetNumber(number int) {
	user.Number = number
}

func (user *User) SetId(id uuid.UUID) {
	user.Id = id
}
