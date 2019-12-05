package receipt

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshalMillistamp(t *testing.T) {

	sampleTime := time.Date(2019, time.March, 12, 10, 11, 12, 0, time.UTC)
	sampleJSON := []byte(`{"value_ms":"1552385472000"}`)

	var data struct {
		Value Millistamp `json:"value_ms,string"`
	}

	if err := json.Unmarshal(sampleJSON, &data); err != nil {
		t.Error(err)
	} else if !sampleTime.Equal(data.Value.Time()) {
		t.Errorf("%v should be the same as %v\n", sampleTime, data.Value.Time())
	}
}
