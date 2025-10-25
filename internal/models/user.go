package models

type User struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Number int    `json:"number"`
}

func (user *User) SetNumber(number int) {
	user.Number = number
}
