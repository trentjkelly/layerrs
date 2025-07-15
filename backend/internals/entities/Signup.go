package entities

type SignupRequest struct {
	Username	string `json:"username"`
	Name		string `json:"name"`
	Email		string `json:"email"`
	Password	string `json:"password"`
}

type SignupResponse struct {

}