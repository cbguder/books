package overdrive

import (
	"encoding/json"
	"time"
)

type FlexibleTime struct {
	time.Time
}

func (t *FlexibleTime) UnmarshalJSON(b []byte) error {
	err := t.Time.UnmarshalJSON(b)
	if err == nil {
		return nil
	}

	var timeString string
	err = json.Unmarshal(b, &timeString)
	if err != nil {
		return err
	}

	t.Time, err = time.Parse("2006-01-02T15:04:05", timeString)
	return err
}
