package receipt

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestGracePeriodExpiration(t *testing.T) {
	data, readErr := ioutil.ReadFile("testdata/response5.json")
	if readErr != nil {
		t.Error(readErr)
	}

	rcpt, parseErr := ParseUnified(data)
	if parseErr != nil {
		t.Error(parseErr)
	}

	graceExpiresAt := time.Date(2020, time.March, 8, 1, 0, 56, 0, time.UTC)
	graceExpiry, ok := rcpt.GracePeriodExpiresDate()
	if ok && !graceExpiry.Equal(graceExpiresAt) {
		t.Errorf("Should parse %s as %s", graceExpiry, graceExpiresAt)
	}

	// if info.Status() != StatusValid {
	// 	t.Error("Should parse status as valid")
	// }
}
