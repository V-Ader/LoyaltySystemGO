package token

type Token struct {
	Id int `json:"id"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
