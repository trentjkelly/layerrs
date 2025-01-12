package entities

import (
	"time"
)

type Artist struct {
	Id 			int			`json:"id"`
	Name 		string		`json:"name"`
	Username 	string		`json:"username"`
	Email 		string		`json:"email"`
	Bio 		string		`json:"bio"`
	R2ImageKey 	string		`json:"r2ImageKey"`
	CreatedAt 	time.Time	`json:"createdAt"`
	UpdatedAt 	time.Time	`json:"updatedAt"`
	Password	string		`json:"password"`
}
