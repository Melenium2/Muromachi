package entities

import (
	"Muromachi/graph/model"
	"fmt"
	"time"
)

// Categories and Keywords representation in db
// Then have same fields
type Track struct {
	Id       int       `json:"-"`
	BundleId int       `json:"bundle,omitempty"`
	Type     string    `json:"type,omitempty"`
	Date     time.Time `json:"date,omitempty"`
	Place    int32     `json:"place,omitempty"`
	App      App       `json:"app,omitempty"`
}

// Converts DBO to *Track or *model.Categories or *model.Keywords
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
