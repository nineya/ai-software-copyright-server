package common

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiredIn    int64  `json:"expired_in"`
	RefreshToken string `json:"refresh_token"`
}
