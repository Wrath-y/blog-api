package util

import (
	"encoding/hex"
	"github.com/satori/go.uuid"
)

func UUID() string {
	uid := hex.EncodeToString(uuid.NewV4().Bytes())
	return uid
}

func UUIDWithLen(length int) string {
	uid := hex.EncodeToString(uuid.NewV4().Bytes())
	return uid[:length]
}
