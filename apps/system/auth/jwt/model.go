package jwt

type Jwt struct {
	// 访问令牌
	AccessToken string `json:"accessToken"`
	// 刷新令牌
	RefreshToken string `json:"refreshToken"`
}
