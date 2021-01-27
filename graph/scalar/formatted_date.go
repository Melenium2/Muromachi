package scalar

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type FormattedDate time.Time

func (date *FormattedDate) UnmarshalGQL(v interface{}) error {
	value, ok := v.(string)
	if !ok {
		return fmt.Errorf("%s", "FormattedDate must be a string")
	}

	timestamp, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}
	*date = FormattedDate(timestamp)

	return nil
}

func (date FormattedDate) MarshalGQL(w io.Writer) {
	_ = binary.Write(w, binary.BigEndian, date)
}

