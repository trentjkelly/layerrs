package entities 

type LoginResponse struct {
	Token string `json:"token"`
	Refresh string `json:"refreshToken"`
}

type JWTResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email 		string 	`json:"email"`
	Password	string 	`json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}