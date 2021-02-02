// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type App struct {
	ID          int       `json:"id"`
	Bundle      string    `json:"bundle"`
	Category    string    `json:"category"`
	DeveloperID string    `json:"developerId"`
	Developer   string    `json:"developer"`
	Geo         string    `json:"geo"`
	StartAt     time.Time `json:"startAt"`
	Period      int       `json:"period"`
}

type Categories struct {
	ID       int       `json:"id"`
	BundleID int       `json:"bundleId"`
	Type     string    `json:"type"`
	Place    int       `json:"place"`
	Date     time.Time `json:"date"`
	App      *App      `json:"app"`
}

type DeveloperContacts struct {
	Email    string `json:"email"`
	Contacts string `json:"contacts"`
}

type Keywords struct {
	ID       int       `json:"id"`
	BundleID int       `json:"bundleId"`
	Type     string    `json:"type"`
	Place    int       `json:"place"`
	Date     time.Time `json:"date"`
	App      *App      `json:"app"`
}

type Meta struct {
	ID               int                `json:"id"`
	BundleID         int                `json:"bundleId"`
	Title            string             `json:"title"`
	Price            string             `json:"price"`
	Picture          string             `json:"picture"`
	Screenshots      []string           `json:"screenshots"`
	Rating           string             `json:"rating"`
	ReviewCount      string             `json:"reviewCount"`
	RatingHistogram  []string           `json:"ratingHistogram"`
	Description      string             `json:"description"`
	ShortDescription string             `json:"shortDescription"`
	RecentChanges    string             `json:"recentChanges"`
	ReleaseDate      string             `json:"releaseDate"`
	LastUpdateDate   string             `json:"lastUpdateDate"`
	Appsize          string             `json:"appsize"`
	Installs         string             `json:"installs"`
	Version          string             `json:"version"`
	OsVersion        string             `json:"osVersion"`
	ContentRating    string             `json:"contentRating"`
	DevContacts      *DeveloperContacts `json:"devContacts"`
	PrivacyPolicy    string             `json:"privacyPolicy"`
	Date             time.Time          `json:"date"`
	App              *App               `json:"app"`
}
