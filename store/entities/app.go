package entities

import (
	"Muromachi/graph/model"
	"fmt"
	"time"
)

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
