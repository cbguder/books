package overdrive

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFlexibleTime(t *testing.T) {
	timeStrings := []string{
		"2024-02-01T12:34:56+00:00",
		"2024-02-01T12:34:56Z",
		"2024-02-01T12:34:56",
	}

	expectedTime := time.Date(2024, 2, 1, 12, 34, 56, 0, time.UTC)

	for _, timeString := range timeStrings {
		jsonString, err := json.Marshal(timeString)
		if err != nil {
			t.Error(err)
		}

		var tt FlexibleTime
		err = tt.UnmarshalJSON(jsonString)
		if err != nil {
			t.Error(err)
		} else if !tt.Equal(expectedTime) {
			t.Errorf("Expected %v, got %v", expectedTime, tt)
		}
	}
}
