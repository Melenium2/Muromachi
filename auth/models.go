package auth

type JWTResponse struct {
	AccessToken  string `json:"access_token,omitempty" form:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty" form:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty" form:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty" form:"refresh_token,omitempty"`
}

type JWTRequest struct {
	AccessType   string `json:"access_type,omitempty" form:"access_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty" form:"refresh_token,omitempty"`
	ClientId     string `json:"client_id,omitempty" form:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty" form:"client_secret,omitempty"`
}
