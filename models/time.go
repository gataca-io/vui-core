package models

import (
	"time"
)

const (
	format = time.RFC3339
)

type TimeWithFormat struct {
	time.Time `swaggertype:"string" json:"tenants" example:"2019-10-01T12:12:05.999Z" description:"Time in RFC 3339 Format"` // embedded time value
	Format    string                                                                                                         `swaggerignore:"true"` // format
}

func (ct TimeWithFormat) MarshalJSON() ([]byte, error) {
	if ct.Format == "" {
		return []byte(`"` + ct.Time.UTC().Format(format) + `"`), nil
	}
	return []byte(`"` + ct.Time.UTC().Format(ct.Format) + `"`), nil
}
