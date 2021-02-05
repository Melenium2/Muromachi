package entities

import (
	"Muromachi/graph/model"
	"fmt"
)

// Type for slice of DBO interface obj
type DboSlice []DBO

// Converts slice of DBO to different types of slices
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
