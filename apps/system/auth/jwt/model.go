package jwt

type Jwt struct {
	// 访问令牌
	AccessToken string `json:"access_token"`
	// 刷新令牌
	RefreshToken string `json:"refresh_token"`
}
