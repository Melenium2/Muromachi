package entities

// Interface for conversion database objects to
// golang struct
type DBO interface {
	To(to interface{}) error
}