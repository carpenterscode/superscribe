package receipt

import (
	"time"
)

type Millistamp int64

func (m Millistamp) Time() time.Time {
	return time.Unix(0, int64(m)*int64(time.Millisecond))
}
