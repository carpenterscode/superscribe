package superscribe

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshalAppleTime(t *testing.T) {

	sampleTime := time.Date(2019, time.March, 12, 10, 11, 12, 0, time.UTC)
	sampleJSON := []byte(`{"SampleTime":"2019-03-12 10:11:12 Etc/GMT"}`)

	var data struct{ SampleTime AppleTime }

	if err := json.Unmarshal(sampleJSON, &data); err != nil {
		t.Error(err)
	} else if !sampleTime.Equal(data.SampleTime.Time) {
		t.Errorf("%v should be the same as %v\n", sampleTime, data.SampleTime.Time)
	}
}
