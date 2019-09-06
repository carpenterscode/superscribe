package receipt

import (
	"encoding/json"
	"time"
)

const appleTimeFormat = "2006-01-02 15:04:05 Etc/GMT"

type AppleTime struct {
	time.Time
}

func (t *AppleTime) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if parsed, err := time.Parse(appleTimeFormat, value); err != nil {
		return err
	} else {
		t.Time = parsed
	}
	return nil
}
