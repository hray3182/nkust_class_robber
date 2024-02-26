package module

type LoginInfo struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ValidateCode string `json:"validateCode"`
}
