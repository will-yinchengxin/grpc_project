package util

import (
	"github.com/go-basic/uuid"
)

func GetUuid() string {
	return uuid.New()
}
