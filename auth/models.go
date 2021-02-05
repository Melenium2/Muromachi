package auth

// Representation of jwt response
type JWTResponse struct {
	// Jwt token for interaction with api
	AccessToken  string `json:"access_token,omitempty" form:"access_token,omitempty"`
	// Token type
	//
	// by default this is 'Bearer'
	TokenType    string `json:"token_type,omitempty" form:"token_type,omitempty"`
	// When the token is expired
	ExpiresIn    int    `json:"expires_in,omitempty" form:"expires_in,omitempty"`
	// Refresh token if in request AccessType was session or refresh_token
	RefreshToken string `json:"refresh_token,omitempty" form:"refresh_token,omitempty"`
}

// Representation of jwt request
type JWTRequest struct {
	// Access type this is type of request.
	//
	// This field should be one of (simple, session, refresh_token)
	AccessType   string `json:"access_type,omitempty" form:"access_type,omitempty"`
	// Token for refreshing
	//
	// Should be with refresh_token AccessType
	RefreshToken string `json:"refresh_token,omitempty" form:"refresh_token,omitempty"`
	// Personal generated client id string
	ClientId     string `json:"client_id,omitempty" form:"client_id,omitempty"`
	// Client secret string
	ClientSecret string `json:"client_secret,omitempty" form:"client_secret,omitempty"`
}
