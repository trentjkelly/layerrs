package entities 

type LoginRequest struct {
	Email 		string 	`json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Refresh string `json:"refreshToken"`
}

type JWTResponse struct {
	Token string `json:"token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}