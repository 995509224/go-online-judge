package help

import uuid "github.com/satori/go.uuid"

func Getuuid() string {
	return uuid.NewV4().String()
}
