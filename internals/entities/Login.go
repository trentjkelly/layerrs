package entities 

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email 		string 	`json:"email"`
	Password	string 	`json:"password"`
}