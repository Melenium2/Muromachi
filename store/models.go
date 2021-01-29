package store

import (
	"Muromachi/graph/model"
	"Muromachi/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jackc/pgtype"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type DBO interface {
	To(to interface{}) error
}

type App struct {
	Id          int       `json:"-"`
	Bundle      string    `json:"bundle,omitempty"`
	Category    string    `json:"category,omitempty"`
	DeveloperId string    `json:"developer_id,omitempty"`
	Developer   string    `json:"developer,omitempty"`
	Geo         string    `json:"geo,omitempty"`
	StartAt     time.Time `json:"start_at,omitempty"`
	Period      uint32    `json:"period,omitempty"`
}

func (a App) To(to interface{}) error {
	switch v := to.(type) {
	case *App:
		*v = a
	case *model.App:
		v.ID = a.Id
		v.Bundle = a.Bundle
		v.Category = a.Category
		v.DeveloperID = a.DeveloperId
		v.Developer = a.Developer
		v.Geo = a.Geo
		v.StartAt = a.StartAt
		v.Period = int(a.Period)
	default:
		return fmt.Errorf("%s", "param 'to' not the same type with *App")
	}

	return nil
}

type Meta struct {
	Id                int               `json:"-"`
	BundleId          int               `json:"-"`
	Title             string            `json:"title" db:"title"`
	Price             string            `json:"price" db:"price"`
	Picture           string            `json:"picture" db:"picture"`
	Screenshots       []string          `json:"screenshots" db:"screenshots"`
	Rating            string            `json:"rating" db:"rating"`
	ReviewCount       string            `json:"reviewCount" db:"review_count"`
	RatingHistogram   []string          `json:"ratingHistogram" db:"rating_histogram"`
	Description       string            `json:"description" db:"description"`
	ShortDescription  string            `json:"shortDescription" db:"short_description"`
	RecentChanges     string            `json:"recentChanges" db:"recent_changes"`
	ReleaseDate       string            `json:"releaseDate" db:"release_date"`
	LastUpdateDate    string            `json:"lastUpdateDate" db:"last_update_date"`
	AppSize           string            `json:"appSize" db:"app_size"`
	Installs          string            `json:"installs" db:"installs"`
	Version           string            `json:"version" db:"version"`
	AndroidVersion    string            `json:"androidVersion" db:"android_version"`
	ContentRating     string            `json:"contentRating" db:"content_rating"`
	DeveloperContacts DeveloperContacts `json:"developerContacts" db:"developer_contacts"`
	PrivacyPolicy     string            `json:"privacyPolicy,omitempty"`
	Date              time.Time         `json:"date,omitempty"`
	App               App               `json:"app,omitempty"`
}

func (m Meta) To(to interface{}) error {
	switch v := to.(type) {
	case *Meta:
		*v = m
	case *model.Meta:
		v.ID = m.Id
		v.BundleID = m.BundleId
		v.Title = m.Title
		v.Price = m.Price
		v.Picture = m.Picture
		v.Screenshots = m.Screenshots
		v.Rating = m.Rating
		v.ReviewCount = m.ReviewCount
		v.RatingHistogram = m.RatingHistogram
		v.Description = m.Description
		v.ShortDescription = m.ShortDescription
		v.RecentChanges = m.RecentChanges
		v.ReleaseDate = m.ReleaseDate
		v.LastUpdateDate = m.LastUpdateDate
		v.Appsize = m.AppSize
		v.Installs = m.Installs
		v.Version = m.Version
		v.OsVersion = m.AndroidVersion
		v.ContentRating = m.ContentRating
		v.DevContacts = &model.DeveloperContacts{
			Email:    m.DeveloperContacts.Email,
			Contacts: m.DeveloperContacts.Email,
		}
		v.PrivacyPolicy = m.PrivacyPolicy
		v.Date = m.Date
		v.App = &model.App{}

		return m.App.To(v.App)
	default:
		return fmt.Errorf("%s", "param 'to' not the same type with *Meta")
	}

	return nil
}

type Track struct {
	Id       int       `json:"-"`
	BundleId int       `json:"bundle,omitempty"`
	Type     string    `json:"type,omitempty"`
	Date     time.Time `json:"date,omitempty"`
	Place    int32     `json:"place,omitempty"`
	App      App       `json:"app,omitempty"`
}

func (tr Track) To(to interface{}) error {
	switch v := to.(type) {
	case *Track:
		*v = tr
	case *model.Categories:
		v.ID = tr.Id
		v.BundleID = tr.BundleId
		v.Type = tr.Type
		v.Date = tr.Date
		v.Place = int(tr.Place)
		v.App = &model.App{}

		return tr.App.To(v.App)
	case *model.Keywords:
		v.ID = tr.Id
		v.BundleID = tr.BundleId
		v.Type = tr.Type
		v.Date = tr.Date
		v.Place = int(tr.Place)
		v.App = &model.App{}

		return tr.App.To(v.App)
	default:
		return fmt.Errorf("%s", "param 'to' not the same type with *Track")
	}

	return nil
}

type DeveloperContacts struct {
	Email    string `json:"email,omitempty"`
	Contacts string `json:"contacts,omitempty"`
}

func (dest *DeveloperContacts) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		return fmt.Errorf("NULL values can't be decoded. Scan into a &*DeveloperContacts to handle NULLs")
	}

	if err := (pgtype.CompositeFields{&dest.Email, &dest.Contacts}).DecodeBinary(ci, src); err != nil {
		return err
	}

	return nil
}

func (src DeveloperContacts) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	email := pgtype.Text{String: src.Email, Status: pgtype.Present}
	contacts := pgtype.Text{String: src.Contacts, Status: pgtype.Present}

	return (pgtype.CompositeFields{&email, &contacts}).EncodeBinary(ci, buf)
}

type DboSlice []DBO

func (d DboSlice) To(to interface{}) error {
	switch v := to.(type) {
	case []*model.App:
		if len(v) != len(d) {
			return fmt.Errorf("len of pointer 'to' not the same with len of DboSlice")
		}
		for i, value := range d {
			app := &model.App{}
			if err := value.To(app); err != nil {
				return err
			}
			v[i] = app
		}
	case []*model.Meta:
		if len(v) != len(d) {
			return fmt.Errorf("len of pointer 'to' not the same with len of DboSlice")
		}
		for i, value := range d {
			meta := &model.Meta{}
			if err := value.To(meta); err != nil {
				return err
			}
			v[i] = meta
		}
	case []*model.Categories:
		if len(v) != len(d) {
			return fmt.Errorf("len of pointer 'to' not the same with len of DboSlice")
		}
		for i, value := range d {
			cat := &model.Categories{}
			if err := value.To(cat); err != nil {
				return err
			}
			v[i] = cat
		}
	case []*model.Keywords:
		if len(v) != len(d) {
			return fmt.Errorf("len of pointer 'to' not the same with len of DboSlice")
		}
		for i, value := range d {
			key := &model.Keywords{}
			if err := value.To(key); err != nil {
				return err
			}
			v[i] = key
		}
	default:
		return fmt.Errorf("param 'to' not the same type with next types ([]*model.App, []*model.Meta, []*model.Categories, []*model.Keywords)")
	}

	return nil
}

type Session struct {
	ID           int       `json:"id,omitempty"`
	UserId       int       `json:"user_id,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	Ip           string    `json:"ip,omitempty"`
	ExpiresIn    time.Time `json:"expires_in,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type User struct {
	ID           int       `json:"id,omitempty"`
	ClientId     string    `json:"client_id,omitempty"`
	ClientSecret string    `json:"client_secret,omitempty"`
	Company      string    `json:"company,omitempty"`
	AddedAt      time.Time `json:"added_at,omitempty"`
}

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

func (u *User) CompareSecret(secret string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.ClientSecret), []byte(secret))
}
