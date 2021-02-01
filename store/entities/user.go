package entities

import (
	"Muromachi/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User representation if db
type User struct {
	ID           int       `json:"id,omitempty"`
	ClientId     string    `json:"client_id,omitempty"`
	ClientSecret string    `json:"client_secret,omitempty"`
	Company      string    `json:"company,omitempty"`
	AddedAt      time.Time `json:"added_at,omitempty"`
}

// Generate random ClientId and ClientSecret for *User struct
func (u *User) GenerateSecrets() error {
	uuid, err := utils.UUID()
	if err != nil {
		return err
	}
	hash := md5.Sum([]byte(uuid))
	u.ClientId = hex.EncodeToString(hash[:])
	u.ClientSecret = utils.Hash(u.ClientId, time.Now().Unix())

	return nil
}

// Func hash clint secret then replace client secret.
//
// Not hashed client secret will return with first return param
func (u *User) SecureSecret() (old string, err error) {
	if u.ClientSecret == "" {
		return "", fmt.Errorf("%s", "empty client secret")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.ClientSecret), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	old = u.ClientSecret
	u.ClientSecret = string(hash)

	return
}

// Compare given secret with userrepo secret
func (u *User) CompareSecret(secret string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.ClientSecret), []byte(secret))
}
