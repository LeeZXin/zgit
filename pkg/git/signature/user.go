package signature

import "time"

type User struct {
	Id          string
	CreatedUnix time.Time
	Name        string
	Email       string
}
