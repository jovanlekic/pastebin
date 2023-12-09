package api


type UserLogin struct {
	Username 	string `json:"username"`
	Password 	string `json:"password"`
}

type UserRegistration struct {
	Username 	string `json:"username"`
	Password 	string `json:"password"`
	Email 	 	string `json:"email"`
}
