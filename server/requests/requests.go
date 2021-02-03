package requests

import "time"

// Request with info about who need to be banned
type TokenList struct {
	// Which user need to be banned
	UserId int      `json:"user_id,omitempty"`
	// Which tokens need to be banned (can be the same with user tokens)
	Tokens []string `json:"tokens,omitempty"`
	// Time for how long to block
	//
	// in seconds
	Ttl    int      `json:"ttl,omitempty"`
}

// Information about ban operations
type BanInfo struct {
	// Type of operation
	//
	// operations: (ban, unban)
	Type   string      `json:"type,omitempty"`
	// How much tokens was banned
	Count  int         `json:"count,omitempty"`
	// Tokens which banned
	Tokens interface{} `json:"tokens,omitempty"`
	// Time
	At     time.Time   `json:"at,omitempty"`
}
