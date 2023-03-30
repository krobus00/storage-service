package utils

import (
	"encoding/base64"
)

func NewBucketKey(key string, id string) string {
	uEnc := base64.URLEncoding.EncodeToString([]byte(id))
	return key + ":" + uEnc[:3]
}
