package entities

import "time"

type Session struct {
	ID           int       `json:"id,omitempty"`
	UserId       int       `json:"user_id,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	Ip           string    `json:"ip,omitempty"`
	ExpiresIn    time.Time `json:"expires_in,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}
