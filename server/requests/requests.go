package requests

import "time"

type TokenList struct {
	UserId int      `json:"user_id,omitempty"`
	Tokens []string `json:"tokens,omitempty"`
	Ttl    int      `json:"ttl,omitempty"`
}

type BanInfo struct {
	Type   string      `json:"type,omitempty"`
	Count  int         `json:"count,omitempty"`
	Tokens interface{} `json:"tokens,omitempty"`
	At     time.Time   `json:"at,omitempty"`
}
